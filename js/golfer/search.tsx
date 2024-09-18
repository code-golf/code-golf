import { $ } from '../_util';

interface Solution {
    Code: string, Hole: string, Lang: string, Scoring: string | null
}
let solutions: Solution[] = JSON.parse($('#solutions').innerText);
solutions = [{"Code":"Solution 1","Hole":"12-days-of-christmas","Lang":"nim","Scoring":"chars"}, {"Code":"Solution 1","Hole":"12-days-of-christmas","Lang":"nim","Scoring":"bytes"},{"Code":"Solution 2","Hole":"12-days-of-christmas","Lang":"python","Scoring":"chars"}, {"Code":"Solution 3 Solution","Hole":"12-days-of-christmas","Lang":"python","Scoring":"bytes"}]
const uniqueSolutions: Record<string, Solution> = {}
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
$('#isRegex').onchange = onSearch;

function onSearch() {
    const search = $<HTMLInputElement>('#searchInput').value;
    if (!search) {
        $('#resultsOverview').innerText = '';
        $<HTMLInputElement>('#searchInput').setCustomValidity('');
        $('#results').replaceChildren();
        return;
    }
    const isRegex = $<HTMLInputElement>('#isRegex').checked;
    let pattern: RegExp | undefined = undefined;
    try {
        pattern = new RegExp(isRegex ? search : search.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), "g");
    }
    catch {
        $<HTMLInputElement>('#searchInput').setCustomValidity('Invalid RegExp');
        $<HTMLInputElement>('#searchInput').reportValidity();
    }
    if (pattern) {
        $<HTMLInputElement>('#searchInput').setCustomValidity('');
        const amount = (n: number, singular: string, plural?: string) => `${n} ${n === 1 ? singular : plural ?? singular + "s"}`;
        const results = solutions.map(x => ({
            ...x,
            matches: [...x.Code.matchAll(pattern)]
        })).filter(x => x.matches.length > 0);
        const totalCount = results.map(x => x.matches.length).reduce((a,b)=>a+b);
        $('#resultsOverview').innerText = results.length === 0
            ? '0 matches'
            : `${amount(totalCount, 'match', 'matches')} across ${amount(results.length, 'solution')}`;
        const resultNodes = results.map(r => (<a href='christmas-trees#python'>
            <h2>{holes[r.Hole]} in {langs[r.Lang]}{r.Scoring ? ` (${r.Scoring})` : ''}</h2>
            <span>{amount(r.matches.length, 'match', 'matches')}</span>
        </a>));
        $('#results').replaceChildren(...resultNodes);
    }
}