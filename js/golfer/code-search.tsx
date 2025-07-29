import { $, debounce } from '../_util';

interface Match {
    before: string, match: string, after: string, count: number, hole: string, lang: string, scoring: string | null
}
const langs: Record<string, string[]> = JSON.parse($('#langs').innerText);
const holes: Record<string, string[]> = JSON.parse($('#holes').innerText);

$('#searchInput').onkeyup = onSearch;
$('#isRegexInput').onchange = onSearch;
$('#languageInput').onchange = onSearch;

$('#languageInput').replaceChildren(<option value=''>All languages</option>, ...Object.entries(langs).map(([id,names]) => <option value={id}>{names[0]}</option>));

const amount = (n: number, singular: string, plural?: string) => `${n} ${n === 1 ? singular : plural ?? singular + 's'}`;

let searchParams = '';

async function onSearch() {
    let pattern = $<HTMLInputElement>('#searchInput').value;
    if (!pattern) {
        searchParams = '';
        $('#resultsOverview').innerText = '';
        $<HTMLInputElement>('#searchInput').setCustomValidity('');
        $('#results').replaceChildren();
        return;
    }
    const isRegexInput = $<HTMLInputElement>('#isRegexInput').checked;
    pattern = isRegexInput ? pattern : pattern.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const lang = $<HTMLSelectElement>('#languageInput').value;

    try {
        new RegExp(pattern);
    }
    catch {
        $<HTMLInputElement>('#searchInput').setCustomValidity('Invalid Regex');
        $<HTMLInputElement>('#searchInput').reportValidity();
        return;
    }

    $<HTMLInputElement>('#searchInput').setCustomValidity('');
    const newSearchParams = new URLSearchParams(lang ? {pattern, lang} : {pattern}).toString();
    if (newSearchParams === searchParams) return;

    searchParams = newSearchParams;
    $('#resultsOverview').innerText = 'searching...';
    fetchSolutionsDebounced();
}

const fetchSolutionsDebounced = debounce(fetchSolutions, 500);

async function fetchSupporting408InFirefox(path: string) {
    try {
        return await fetch(path);
    }
    catch (e) {
        if (`${e}`.includes('NetworkError')) return {status: 408} as const;
        throw e;
    }
}

async function fetchSolutions() {
    const resp = await fetchSupporting408InFirefox(`/api/solutions-search?${searchParams}`);
    if (resp.status !== 200) {
        $<HTMLInputElement>('#searchInput').setCustomValidity(resp.status === 408 ? 'Timeout. Try single language or raw text search.' : `Error ${resp.status}`);
        $<HTMLInputElement>('#searchInput').reportValidity();
        return;
    }
    const results = await resp.json() as Match[];
    const totalCount = results.map(x => x.count).reduce((a,b)=>a+b, 0);
    const holesCount = [...new Set(results.map(x => x.hole))].length;
    const resultsLangs = [...new Set(results.map(x => x.lang))].map(x => langs[x]);
    $('#resultsOverview').innerText = results.length === 0
        ? '0 matches'
        : `${amount(totalCount, 'match', 'matches')} across ${amount(results.length, 'solution')} (${amount(holesCount, 'hole')} in ${resultsLangs.length > 5 ? `${resultsLangs.length} languages` : resultsLangs.join(', ')})`;
    const resultNodes = results.map(r => (<a href={'/' + r.hole + '#' + r.lang}>
        <h2>{holes[r.hole]?.[0]} in {langs[r.lang]?.[0]}{r.scoring ? ` (${r.scoring})` : ''}</h2>
        <span>
            <code>{r.before.split(/\n|\r/).at(-1)}<span class='match'>{r.match}</span>{r.after.split(/\n|\r/)[0]}</code>
            {r.count > 1 ? ' and ' + amount(r.count-1, 'more match', 'more matches') : ''}
        </span>
    </a>));
    $('#results').replaceChildren(...resultNodes);
}
