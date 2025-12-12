import { $, amount, debounce } from '../_util';

interface Match {
    before: string, match: string, after: string, count: number, hole: string, lang: string, scoring: string | null
}
const langs: Record<string, string[]> = JSON.parse($('#langs').innerText);
const holes: Record<string, string[]> = JSON.parse($('#holes').innerText);

const form = $<HTMLFormElement>('#search');

let searchParams = '';

form.onsubmit = e => e.preventDefault();

form.onchange = form.q.onkeyup = onload = async () => {
    const hole    = form.hole.value;
    const lang    = (form.lang as any).value;
    const pattern = form.regex.checked
        ? form.q.value : form.q.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

    try {
        new RegExp(pattern);
    }
    catch {
        form.q.setCustomValidity('Invalid Regex');
        form.q.reportValidity();
        return;
    }

    form.q.setCustomValidity('');

    const newSearchParams = new URLSearchParams({pattern, hole, lang}).toString();
    if (newSearchParams === searchParams) return;
    searchParams = newSearchParams;

    if (pattern != '') {
        $('#resultsOverview').innerText = 'searching...';
        fetchSolutionsDebounced();
    }
    else {
        $('#resultsOverview').innerText = $('#results').innerHTML = '';
    }
};

const fetchSolutionsDebounced = debounce(fetchSolutions, 500);

async function fetchSolutions() {
    const resp = await fetch(`/api/solutions-search?${searchParams}`);
    if (resp.status !== 200) {
        form.q.setCustomValidity(resp.status === 504 ?
            'Timeout. Try a less complex search' : `Error ${resp.status}`);
        form.q.reportValidity();
        return;
    }
    const results = await resp.json() as Match[];
    const totalCount = results.map(x => x.count).reduce((a,b)=>a+b, 0);
    const holesCount = [...new Set(results.map(x => x.hole))].length;
    const resultsLangs = [...new Set(results.map(x => x.lang))].map(x => langs[x][0]);
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
