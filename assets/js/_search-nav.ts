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
    function normalize(str: string): string {
        return str.toLowerCase()
            .replace(/[áàäâãå]/g, 'a')
            .replace(/[éèëê]/g, 'e')
            .replace(/[íìïî]/g, 'i')
            .replace(/[óòöôõ]/g, 'o')
            .replace(/[úùüû]/g, 'u')
            .replace(/[ýÿ]/g, 'y')
            .replace(/[č]/g, 'c')
            .replace(/[ď]/g, 'd')
            .replace(/[ě]/g, 'e')
            .replace(/[ň]/g, 'n')
            .replace(/[ř]/g, 'r')
            .replace(/[š]/g, 's')
            .replace(/[ť]/g, 't')
            .replace(/[ž]/g, 'z');
    }

    needle = normalize(needle);
    haystack = normalize(haystack);
    if (needle === haystack) return 1e11;

    const isJustLetters = /^[a-z]+$/.test(needle);
    if (haystack.startsWith(needle)) return 1e9 - haystack.length;
    if (isJustLetters) {
        const fuzzyRegex = new RegExp('([^a-z]|^)' + needle.split('').join('(.*([^a-z]|^))?'));
        if (fuzzyRegex.test(haystack)) return 1e7 - haystack.split(/[^a-z]/).length*100 - haystack.length;
    }
    if (haystack.includes(needle)) return 1e5 - haystack.length - needle.length;
    if (isJustLetters) {
        const fuzzyRegex = new RegExp(needle.split('').join('.*'));
        if (fuzzyRegex.test(haystack)) return 1e3 - haystack.length;
    }
    return 0;
}

let lastSearch: string | undefined = undefined;
// eslint-disable-next-line no-unused-vars
export function requestResults(search: string, updateResults: (results: SearchNavResult[]) => void) {
    if (lastSearch === search) return;
    lastSearch = search;

    if (!search) {
        updateResults([]);
        return;
    }

    if (search.startsWith('@')) {
        requestAtResults(search, updateResults);
        return;
    }

    const [holeSearch, langSearch] = search.includes('#')
        ? [search.slice(0, search.indexOf('#')), search.slice(search.indexOf('#') + 1)]
        : [search, undefined];
    const currentHole = location.pathname.slice(1) in holes ? location.pathname.slice(1) : undefined;

    let matches: SearchNavResultInternal[] = [];

    const holesMatches = Object.entries(holes)
        .map(([k,v]) => ({
            path: `/${k}`,
            description: v[0],
            priority: Math.max(...v.map(name => defaultMatchPriority(name, holeSearch))),
        })).filter(x => x.priority > 0);

    const langsMatches = Object.entries(langs)
        .map(([k,v]) => ({
            path: `#${k}`,
            description: v[0],
            priority: Math.max(...v.map(name => defaultMatchPriority(name, langSearch ?? holeSearch) * (langSearch ? 1 : 0.1))),
        })).filter(x => x.priority > 0);

    if (langSearch && holeSearch) {
        matches = holesMatches.flatMap(hole => langsMatches.map(lang => ({
            path: hole.path + lang.path,
            description: hole.description + ' in ' + lang.description,
            priority: hole.priority * lang.priority,
        })));
    }
    else if (holeSearch) {
        matches = holesMatches;
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

    updateResults(processResults(matches));
}

function processResults(results: SearchNavResultInternal[]) {
    results = results.filter(x => x.priority > 0 && new URL(x.path, location.origin).href !== location.href);
    const maxPriority = Math.max(...results.map(x => x.priority), 0);
    results = results.filter(x => x.priority >= maxPriority / 1000);
    results.sort((a,b) => b.priority - a.priority);
    return results.slice(0, 10);
}

// eslint-disable-next-line no-unused-vars
function requestAtResults(search: string, updateResults: (results: SearchNavResult[]) => void) {
    if (search === '@') {
        const currentGolferPath = $<HTMLAnchorElement>('#site-header [title=Profile]')?.href;
        updateResults(currentGolferPath ? processResults([{path: new URL(currentGolferPath).pathname, description: 'My Profile', priority: 1}]) : []);
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
        path: `/golfers/${x}`,
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