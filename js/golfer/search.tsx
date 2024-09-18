import { $ } from '../_util';

interface Solution {
    Code: string, Hole: string, Lang: string, Scoring: string | null
}
let solutions: Solution[] = JSON.parse($('#solutions').innerText);
const uniqueSolutions: Record<string, Solution> = {};
for (const solution of solutions) {
    const key = solution.Hole + solution.Lang + solution.Code;
    if (key in uniqueSolutions) {
        uniqueSolutions[key].Scoring = null;
    }
    else {
        uniqueSolutions[key] = solution;
    }
}
solutions = Object.values(uniqueSolutions);
const langs: Record<string,string> = JSON.parse($('#langs').innerText);
const holes: Record<string,string> = JSON.parse($('#holes').innerText);

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
                .filter(x => !language || x.Lang == language)
                .map(x => {
                    const matches = [...x.Code.matchAll(pattern)];
                    const matchesCount = matches.length;
                    let firstMatch = {before: '', match: '', after: ''};
                    if (matchesCount > 0) {
                        const m = matches[0];
                        const b = m.index;
                        let a = b;
                        while (a-1 > 0 && a > b - 10 && !/\n|\r/.test(x.Code[a-1])) a--;
                        let c = b;
                        while (c+1 <= b + m[0].length && !/\n|\r/.test(x.Code[c])) c++;
                        let d = c;
                        while (d+1 <= x.Code.length && d < c + 10 && !/\n|\r/.test(x.Code[d])) d++;
                        firstMatch = {before: x.Code.slice(a,b), match: x.Code.slice(b,c), after: x.Code.slice(c,d)};
                    }
                    return {...x, matchesCount, firstMatch};
                }).filter(x => x.matchesCount > 0);
            results.sort((a,b)=> a.matchesCount - b.matchesCount);
            const totalCount = results.map(x => x.matchesCount).reduce((a,b)=>a+b, 0);
            const ci = caseInsensitive ? 'case insensitive ' : '';
            $('#resultsOverview').innerText = results.length === 0
                ? '0 matches'
                : `${amount(totalCount, ci + 'match', ci + 'matches')} across ${amount(results.length, 'solution')}`;
            const resultNodes = results.map(r => (<a href={r.Hole + '#' + r.Lang}>
                <h2>{holes[r.Hole]} in {langs[r.Lang]}{r.Scoring ? ` (${r.Scoring})` : ''}</h2>
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