import { $ } from '../_util';

interface Solution {
    code: string, hole: string, lang: string, scoring: string | null
}
let solutions: Solution[] = JSON.parse($('#solutions').innerText);
const uniqueSolutions: Record<string, Solution> = {};
for (const solution of solutions) {
    const key = solution.hole + solution.lang + solution.code;
    if (key in uniqueSolutions) {
        uniqueSolutions[key].scoring = null;
    }
    else {
        uniqueSolutions[key] = solution;
    }
}
solutions = Object.values(uniqueSolutions);
const langs: Record<string, string> = JSON.parse($('#langs').innerText);
const holes: Record<string, string> = JSON.parse($('#holes').innerText);

$('#searchInput').onkeyup = onSearch;
$('#isRegexInput').onchange = onSearch;
$('#languageInput').onchange = onSearch;

$('#languageInput').replaceChildren(<option value=''>All languages</option>, ...Object.entries(langs).map(([id,name]) => <option value={id}>{name}</option>));

function onSearch() {
    let search = $<HTMLInputElement>('#searchInput').value;
    if (!search) {
        $('#resultsOverview').innerText = '';
        $<HTMLInputElement>('#searchInput').setCustomValidity('');
        $('#results').replaceChildren();
        return;
    }
    const isRegexInput = $<HTMLInputElement>('#isRegexInput').checked;
    search = isRegexInput ? search : search.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const language = $<HTMLSelectElement>('#languageInput').value;

    // If there is 0 case sensitive matches, we match case-insensitively
    for (const caseInsensitive of [false, true]) {
        let pattern: RegExp | undefined = undefined;
        try {
            pattern = new RegExp(search, caseInsensitive ? 'gi' : 'g');
        }
        catch {
            $<HTMLInputElement>('#searchInput').setCustomValidity('Invalid Regex');
            $<HTMLInputElement>('#searchInput').reportValidity();
        }
        if (pattern) {
            $<HTMLInputElement>('#searchInput').setCustomValidity('');
            const amount = (n: number, singular: string, plural?: string) => `${n} ${n === 1 ? singular : plural ?? singular + 's'}`;

            const results = solutions
                .filter(x => !language || x.lang == language)
                .map(x => {
                    const matches = [...x.code.matchAll(pattern)];
                    const matchesCount = matches.length;
                    let firstMatch = {before: '', match: '', after: ''};
                    if (matchesCount > 0) {
                        const m = matches[0];
                        const b = m.index;
                        let a = b;
                        while (a-1 > 0 && a > b - 10 && !/\n|\r/.test(x.code[a-1])) a--;
                        let c = b;
                        while (c+1 <= b + m[0].length && !/\n|\r/.test(x.code[c])) c++;
                        let d = c;
                        while (d+1 <= x.code.length && d < c + 10 && !/\n|\r/.test(x.code[d])) d++;
                        firstMatch = {before: x.code.slice(a,b), match: x.code.slice(b,c), after: x.code.slice(c,d)};
                    }
                    return {...x, matchesCount, firstMatch};
                }).filter(x => x.matchesCount > 0);
            results.sort((a,b)=> a.matchesCount - b.matchesCount);
            const totalCount = results.map(x => x.matchesCount).reduce((a,b)=>a+b, 0);
            const ci = caseInsensitive ? 'case insensitive ' : '';
            $('#resultsOverview').innerText = results.length === 0
                ? '0 matches'
                : `${amount(totalCount, ci + 'match', ci + 'matches')} across ${amount(results.length, 'solution')}`;
            const resultNodes = results.map(r => (<a href={r.hole + '#' + r.lang}>
                <h2>{holes[r.hole]} in {langs[r.lang]}{r.scoring ? ` (${r.scoring})` : ''}</h2>
                <span>
                    <code>{r.firstMatch.before}<span class='match'>{r.firstMatch.match}</span>{r.firstMatch.after}</code>
                    {r.matchesCount > 1 ? ' and ' + amount(r.matchesCount-1, 'more match', 'more matches') : ''}
                </span>
            </a>));
            $('#results').replaceChildren(...resultNodes);
            if (totalCount > 0) break;
        }
    }
}
