import { $ } from './_util';

export interface SearchNavResult {
    path: string;
    description: string;
}

interface SearchNavResultInternal extends SearchNavResult {
    priority: number;
}

const langs: Record<string, string[]> = JSON.parse($('#langs').innerText);
const holes: Record<string, string[]> = JSON.parse($('#holes').innerText);

function defaultMatchPriority(haystack: string, needle: string) {
    needle = needle.toLowerCase();
    haystack = haystack.toLowerCase();
    if (needle === haystack) return 1e10;
    if (haystack.startsWith(needle)) return 1e5 - haystack.length;
    if (haystack.includes(needle)) return 1e3 - haystack.length - needle.length;
    return 0;
}

let lastSearch: string | undefined = undefined;
// eslint-disable-next-line no-unused-vars
export function requestResults(search: string, updateResults: (results: SearchNavResult[]) => void) {
    if (!search || lastSearch === search) return;
    lastSearch = search;

    if (search.startsWith('@')) {
        requestAtResults(search, updateResults);
        return;
    }

    const [holeSearch, langSearch] = search.includes('#') ? search.split('#', 2) : [search, undefined];
    const currentHole = location.pathname.slice(1) in holes ? location.pathname.slice(1) : undefined;

    let matches: SearchNavResultInternal[] = [];

    const holesMatches = Object.entries(holes)
        .map(([k,v]) => ({
            path: `/${k}`,
            description: v[0],
            priority: Math.max(...v.map(name => defaultMatchPriority(name, holeSearch))),
        })).filter(x => x.priority > 0);

    if (langSearch) {
        const langsMatches = langSearch ? Object.entries(langs)
            .map(([k,v]) => ({
                path: `#${k}`,
                description: v[0],
                priority: Math.max(...v.map(name => defaultMatchPriority(name, langSearch))),
            })).filter(x => x.priority > 0) : [];

        if (holeSearch) {
            matches = holesMatches.flatMap(hole => langsMatches.map(lang => ({
                path: hole.path + lang.path,
                description: hole.description + ' in ' + lang.description,
                priority: hole.priority * lang.priority
            })));
        }
        if (currentHole) {
            matches.push(...langsMatches.map(lang => (
                {
                    path: '/' + currentHole + lang.path,
                    description: holes[currentHole][0] + ' in ' + lang.description,
                    priority: lang.priority,
                }
            )));
        }
    }
    else {
        matches = holesMatches;
    }

    updateResults(processResults(matches));
}

function processResults(results: SearchNavResultInternal[]) {
    results = results.filter(x => x.priority > 0 && new URL(x.path, location.origin).href !== location.href);
    results.sort((a,b) => b.priority - a.priority);
    return results.slice(0, 10);
}

// eslint-disable-next-line no-unused-vars
function requestAtResults(search: string, updateResults: (results: SearchNavResult[]) => void) {
    if (search === '@') {
        const currentGolferPath = $<HTMLAnchorElement>('#site-header [title=Profile]')?.href;
        updateResults(currentGolferPath ? processResults([{path: new URL(currentGolferPath).href, description: 'My Profile', priority: 0}]) : []);
        return;
    }
    requestGolferResults(search.slice(1), updateResults);
}

const requestGolferResults = debounce(fetchGolferResults);

// eslint-disable-next-line no-unused-vars
async function fetchGolferResults(search: string, updateResults: (results: SearchNavResult[]) => void) {
    const resp  = await fetch(`/api/suggestions/golfers?${new URLSearchParams({q: search})}`);
    const golfers = await resp.json() as string[];
    updateResults(processResults(golfers.map(x => ({
        path: `golfers/${x}`,
        description: x,
        priority: defaultMatchPriority(x, search),
    }))));
}

function debounce(func: Function, timeout = 500) {
    let timer: number | undefined;
    return (...args: any[]) => {
        clearTimeout(timer);
        timer = setTimeout(() => func(...args), timeout);
    };
}