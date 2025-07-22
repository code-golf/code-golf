import { ASMStateField }                       from '@defasm/codemirror';
import { $, $$, byteLen, charLen, comma, ord } from './_util';
import { Vim }                                 from '@replit/codemirror-vim';
import { EditorState, EditorView, extensions } from './_codemirror';
import LZString                                from 'lz-string';

let tabLayout: boolean = false;

const langWikiCache: Record<string, string | null> = {};
async function getLangWikiContent(lang: string): Promise<string> {
    if (!(lang in langWikiCache)) {
        langWikiCache[lang] = await _getLangWikiContent(lang);
    }
    return langWikiCache[lang] ?? 'No data for current lang.';
}
async function _getLangWikiContent(lang: string): Promise<string | null> {
    const resp  = await fetch(`/api/wiki/langs/${lang}`);
    const {content, title} = await resp.json() as {content: string, title: string};
    if (resp.status !== 200) {
        return null;
    }
    const header = (<p id={`lang-wiki-${lang}`}>
        Wiki: {title}{' '}
        <a href={`/wiki/langs/${lang}`} target="_blank">
            (open in new tab)
        </a>
        .
    </p>).outerHTML;
    return header + content;
}

const holeLangNotesCache: Record<string, string | null> = {};
async function getHoleLangNotesContent(lang: string): Promise<string> {
    if (!(lang in holeLangNotesCache)) {
        const resp  = await fetch(`/api/notes/${hole}/${lang}`);
        holeLangNotesCache[lang] = resp.status === 200 ? (await resp.text()) : null;
    }
    return holeLangNotesCache[lang] ?? '';
}

const renamedHoles: Record<string, string> = {
    'billiard':                      'billiards',
    'eight-queens':                  'n-queens',
    'factorial-factorisation-ascii': 'factorial-factorisation',
    'grid-packing':                  'css-grid',
    'placeholder':                   'tutorial',
    'sudoku-v2':                     'sudoku-fill-in',
};

const renamedLangs: Record<string, string> = {
    lisp:  'common-lisp',
    perl6: 'raku',
};

export function init(_tabLayout: boolean, setSolution: any, setCodeForLangAndSolution: any, updateReadonlyPanels: any, getEditor: () => any) {
    tabLayout = _tabLayout;
    const closuredSubmit = () => submit(getEditor(), updateReadonlyPanels);
    if (vimMode) Vim.defineEx('write', 'w', closuredSubmit);

    (onhashchange = async () => {
        const hashLang = location.hash.slice(1) || localStorage.getItem('lang');

        // Kick 'em to Python if we don't know the chosen language, or if there is no given language.
        lang = hashLang && langs[hashLang] ? hashLang : 'python';

        // Assembly only has bytes.
        if (lang == 'assembly')
            setSolution(0);

        localStorage.setItem('lang', lang);

        history.replaceState(null, '', '#' + lang);

        const editor = getEditor();
        if (tabLayout) refreshScores(editor);
        setCodeForLangAndSolution(editor);

        if (tabLayout) {
            updateReadonlyPanels({langWiki: await getLangWikiContent(lang)});
            updateReadonlyPanels({holeLangNotes: await getHoleLangNotesContent(lang)});
        }
    })();

    $('dialog [name=text]').addEventListener('input', (e: Event) => {
        const target = e.target as HTMLInputElement;
        target.form!.confirm.toggleAttribute('disabled',
            target.value !== target.placeholder);
    });

    for (const [key, value] of Object.entries(localStorage)) {
        if (key.startsWith('code_')) {
            const [prefix, hole, lang, scoring] = key.split('_');

            const newHole = renamedHoles[hole] ?? hole;
            const newLang = renamedLangs[lang] ?? lang;

            const newKey = [prefix, newHole, newLang, scoring].join('_');
            if (key !== newKey) {
                localStorage.setItem(newKey, value);
                localStorage.removeItem(key);
            }
        }
    }
}

export function initDeleteBtn(deleteBtn: HTMLElement | undefined, langs: any) {
    deleteBtn?.addEventListener('click', () => {
        $('#delete-dialog b').innerText = langs[lang].name;
        $<HTMLInputElement>('#delete-dialog [name=lang]').value = lang;
        $<HTMLInputElement>('#delete-dialog [name=text]').value = '';
        $<HTMLDialogElement>('#delete-dialog').showModal();
    });
}

export function initCopyButtons(buttons: NodeListOf<HTMLElement>) {
    for (const btn of buttons)
        btn.onclick = () =>
            navigator.clipboard.writeText(btn.dataset.copy!);
}

export const langs = JSON.parse($('#langs').innerText);
const sortedLangs  =
    Object.values(langs).sort((a: any, b: any) => a.name.localeCompare(b.name));
let lang: string = '';

export function getLang() {
    return lang;
}

export const hole         = decodeURI(location.pathname.slice(1));
const scorings     = ['Bytes', 'Chars'];
const solutions    = JSON.parse($('#solutions').innerText);

const vimMode = JSON.parse($('#keymap').innerText) === 'vim';
const vimModeExtensions = vimMode ? [extensions.vim] : [];

const baseExtensions = [...vimModeExtensions, ...extensions.base, ...extensions.editor];

let latestSubmissionID = 0;
let solution = scorings.indexOf(localStorage.getItem('solution') ?? 'Bytes') as 0 | 1;
let scoring  = scorings.indexOf(localStorage.getItem('scoring')  ?? 'Bytes') as 0 | 1;

let hideDeleteBtn: boolean = false;

// The savedInDB state is used to avoid saving solutions in localStorage when
// those solutions match the solutions in the database. It's used to avoid
// restoring a solution from localStorage when the user has improved that
// solution on a different browser. Assume the user is logged-in by default.
// At this point, it doesn't matter whether the
// user is actually logged-in, because solutions dictionaries will be empty
// for users who aren't logged-in, so the savedInDB state won't be used.
// By the time they are non-empty, the savedInDB state will have been updated.
let savedInDB = true;

export function getSavedInDB() {
    return savedInDB;
}

export function getAutoSaveKey(lang: string, solution: 0 | 1) {
    return `code_${hole}_${lang}_${solution}`;
}

function getOtherScoring(value: 0 | 1) {
    return 1 - value as 0 | 1;
}

function getScoring(str: string, index: 0 | 1) {
    return scorings[index] == 'Bytes' ? byteLen(str) : charLen(str);
}

function getSolutionCode(lang: string, solution: 0 | 1): string {
    return lang in solutions[solution] ? solutions[solution][lang] : '';
}

/**
 * Get the code corresponding to the current lang and solution (bytes/chars)
 */
export function getCurrentSolutionCode(): string {
    return getSolutionCode(lang, solution);
}

export function setSolution(value: 0 | 1) {
    solution = value;
    localStorage.setItem('solution', scorings[solution]);
}

export function getSolution() {
    return solution;
}

export function setState(code: string, editor: EditorView) {
    if (!editor) return;
    editor.setState(
        EditorState.create({
            doc: code,
            extensions: [
                ...baseExtensions,
                extensions[lang as keyof typeof extensions] ?? [],
                // These languages shouldn't match brackets.
                ['fish', 'hexagony'].includes(lang)
                    ? extensions.zeroIndexedLineNumbers : [extensions.lineNumbers, extensions.bracketMatching],
                // These languages shouldn't wrap lines.
                ['assembly', 'fish', 'hexagony'].includes(lang)
                    ? [] : EditorView.lineWrapping,
            ],
        }),
    );
    editor.dispatch();  // Dispatch to update strokes.
}

export function setCode(code: string, editor: EditorView | null) {
    editor?.dispatch({
        changes: { from: 0, to: editor.state.doc.length, insert: code },
    });
}

function updateLangPicker() {
    const selectNodes: Node[] = [];
    const langSelect = <select><option value="">Other</option></select>;
    const experimentalLangGroup = <optgroup label="Experimental"></optgroup>;
    let currentLangUnused = false;

    for (const l of sortedLangs as any[]) {
        if (!getSolutionCode(l.id, 0) &&
            !localStorage.getItem(getAutoSaveKey(l.id, 0)) &&
            !localStorage.getItem(getAutoSaveKey(l.id, 1))) {
            const parent = l.experiment ? experimentalLangGroup : langSelect;
            parent.appendChild(<option value={l.id}>{l.name}</option>);
            currentLangUnused ||= lang == l.id;
        }
    }

    if (experimentalLangGroup.childElementCount > 0)
        langSelect.appendChild(experimentalLangGroup);

    if (langSelect.childElementCount > 1) {
        langSelect.addEventListener('change', (e: Event) => {
            const target = e.target as HTMLSelectElement;
            location.hash = '#' + target.value;
        });

        langSelect.value = currentLangUnused ? lang : '';
        if (currentLangUnused) {
            langSelect.classList.add('selectActive');
        }
        selectNodes.push(langSelect);
    }

    // Hybrid language selector: make it easy to see your existing solutions and their lengths.
    const picker = $('#picker');
    const icon   = picker.dataset.style?.includes('icon')  ?? true;
    const label  = picker.dataset.style?.includes('label') ?? true;
    picker.replaceChildren(...sortedLangs.map((l: any) => {
        const tab = <a href={l.id == lang ? null : '#'+l.id} title={l.name}></a>;

        if (icon)  tab.append(<svg><use href={l['logo-url']+'#a'}/></svg>);
        if (label) tab.append(l.name);

        if (getSolutionCode(l.id, 0)) {
            const bytes = byteLen(getSolutionCode(l.id, 0));
            const chars = charLen(getSolutionCode(l.id, 1));

            let text = comma(bytes);
            if (chars && bytes != chars) text += '/' + comma(chars);

            tab.append(<sup>{text}</sup>);
        }
        else if (!localStorage.getItem(getAutoSaveKey(l.id, 0)) &&
                 !localStorage.getItem(getAutoSaveKey(l.id, 1))) {
            return null;
        }

        if (l.experiment) tab.append(<svg><use href="#flask"/></svg>);

        return tab;
    }).filter((x: Node | null) => x), ...selectNodes);
}

export async function refreshScores(editor: any) {
    updateLangPicker();

    // Populate (and show) the solution picker if necessary.
    //
    // We have two database solutions (or local solutions) and they differ.
    // Or if a logged-in user has an auto-saved solution for the other metric,
    // that they have not submitted since logging in, they must be allowed to
    // switch to it, so they can submit it.
    const dbBytes = getSolutionCode(lang, 0);
    const dbChars = getSolutionCode(lang, 1);
    const lsBytes = localStorage.getItem(getAutoSaveKey(lang, 0));
    const lsChars = localStorage.getItem(getAutoSaveKey(lang, 1));

    if ((dbBytes && dbChars && dbBytes != dbChars)
     || (lsBytes && lsChars && lsBytes != lsChars)
     || (dbBytes && lsChars && dbBytes != lsChars && solution == 0)
     || (lsBytes && dbChars && lsBytes != dbChars && solution == 1)) {
        $('#solutionPicker').replaceChildren(...scorings.map((scoring, iNumber) => {
            const i = iNumber as 0 | 1;
            const a = <a>Fewest {scoring}</a>;

            const code = getSolutionCode(lang, i);
            if (code) a.append(' ', <sup>{comma(getScoring(code, i))}</sup>);

            if (i != solution) {
                a.href = '';
                a.onclick = (e: MouseEvent) => {
                    e.preventDefault();
                    setSolution(i);
                    setCodeForLangAndSolution(editor);
                };
            }

            return a;
        }));

        $('#solutionPicker').classList.remove('hide');
    }
    else
        $('#solutionPicker').classList.add('hide');

    // Hide the delete button if we have no solutions.
    hideDeleteBtn = !dbBytes && !dbChars;
    $('#deleteBtn')?.classList.toggle('hide', hideDeleteBtn);

    await populateScores(editor);
}

export function getHideDeleteBtn() {
    return hideDeleteBtn;
}

function tooltip(row: any, scoring: 'Bytes' | 'Chars') {
    const bytes = scoring === 'Bytes' ? row.bytes : row.chars_bytes;
    const chars = scoring === 'Chars' ? row.chars : row.bytes_chars;

    if (bytes === null) return;

    return `${scoring} solution is ${comma(bytes)} bytes` +
        (chars !== null ? `, ${comma(chars)} chars.` : '.');
}

export interface RankFromTo {
    joint: boolean | null,
    rank: number | null,
    strokes: number | null,
}

export interface RankUpdate {
    scoring: string,
    failingStrokes: number | null,
    from: RankFromTo,
    to: RankFromTo,
    oldBestCurrentGolferCount: number | null,
    oldBestCurrentGolferId: number | null,
    oldBestFirstGolferID: number | null,
    oldBestStrokes: number | null,
    oldBestSubmitted: string | null,
    newSolutionCount: number | null,
}

export interface Run {
    answer: string,
    output_delimiter: string,
    multiset_item_delimiter: string,
    args: string[],
    exit_code: number,
    pass: boolean,
    stderr: string,
    stdout: string,
    time_ns: number,
    timeout: boolean,
}

export interface ReadonlyPanelsData {
    Pass: boolean,
    Out: string,
    Exp: string,
    Err: string,
    Argv: string[],
    OutputDelimiter: string,
    MultisetItemDelimiter: string
}

export interface SubmitResponse {
    cheevos:      { emoji: string, name: string }[],
    logged_in:    boolean,
    rank_updates: RankUpdate[],
    runs:         Run[]
}

const makeSingular = (strokes: number, units: string) =>
    strokes == 1 ? units.substring(0, units.length - 1) : units;

const getDisplayRank = (rank: number) => rank < 4 ? ['🥇','🥈','🥉'][rank - 1] : `${rank}${ord(rank)} place`;

// Don't show the delta, if it's the first time playing this hole.
const getDisplayRankChange = (rank: number, delta: number) =>
    getDisplayRank(rank) + (delta > 0 ? ` (was ${getDisplayRank(rank + delta)})` : '');

const getDisplayPointsChange = (points: number, delta: number) =>
    `${points} points` + (delta > 0 && delta < points ? ` (+${delta})` : '');

const scorePopups = (updates: RankUpdate[]) => {
    const strokesDelta = [0, 0];
    const pointsDelta = [0, 0];
    const points = [0, 0];
    const rankDelta = [0, 0];
    const rank = [0, 0];
    let newSolution = false;

    for (const [i, update] of updates.entries()) {
        if (update.to.strokes) {
            const newBest = update.oldBestStrokes != null ?
                Math.min(update.oldBestStrokes, update.to.strokes) :
                update.to.strokes;
            points[i] = Math.round(newBest / update.to.strokes * 1000);

            if (update.from.strokes) {
                strokesDelta[i] = update.from.strokes - update.to.strokes;
                if (update.oldBestStrokes) {
                    pointsDelta[i] = points[i] - Math.round(update.oldBestStrokes / update.from.strokes * 1000);
                }
            }
            else {
                newSolution = true;
                pointsDelta[i] = points[i];
            }
        }

        if (update.to.rank) {
            rank[i] = update.to.rank;
            rankDelta[i] = (update.from.rank || 0) - update.to.rank;
        }
    }

    const nodes: Node[] = [];

    if (strokesDelta[0] > 0 || strokesDelta[1] > 0) {
        let amount = '';

        // Show the decrease in strokes.
        if (strokesDelta[0] > 0 && strokesDelta[0] == strokesDelta[1]) {
            const delta = strokesDelta[0];
            let units = '';
            for (const i of [0, 1] as const) {
                units += (i == 1 ? '/' : '') + makeSingular(delta, updates[i].scoring);
            }

            amount = `${delta} ${units}`;
        }
        else {
            for (const i of [0, 1] as const) {
                if (strokesDelta[i] > 0) {
                    amount += (i == 1 && strokesDelta[0] > 0 ? '/' : '') + `${strokesDelta[i]} ${makeSingular(strokesDelta[i], updates[i].scoring)}`;
                }
            }
        }

        nodes.push(<h3>Score Improved</h3>);
        nodes.push(<p>Saved {amount}</p>);
    }
    else if (newSolution) {
        nodes.push(<h3>New Solution</h3>);
    }

    // Show points updates, including the current number of points, because this is not show on the mini-scoreboard.
    if (pointsDelta[0] != 0 && pointsDelta[0] == pointsDelta[1] && points[0] == points[1]) {
        nodes.push(<p>{getDisplayPointsChange(points[0], pointsDelta[0])} for {updates[0].scoring}/{updates[1].scoring}</p>);
    }
    else {
        for (const i of [0, 1] as const) {
            if (pointsDelta[i] != 0) {
                nodes.push(<p>{getDisplayPointsChange(points[i], pointsDelta[i])} for {updates[i].scoring}</p>);
            }
        }
    }

    // Show rank update.
    if (rankDelta[0] != 0 && rankDelta[0] == rankDelta[1] && rank[0] == rank[1]) {
        nodes.push(<p>{getDisplayRankChange(rank[0], rankDelta[0])} for {updates[0].scoring}/{updates[1].scoring}</p>);
    }
    else {
        for (const i of [0, 1] as const) {
            if (rankDelta[i] != 0) {
                nodes.push(<p>{getDisplayRankChange(rank[i], rankDelta[i])} for {updates[i].scoring}</p>);
            }
        }
    }

    return nodes.length > 1 ? [<div>{nodes}</div>] : [];
};

const diamondPopups = (updates: RankUpdate[]) => {
    const popups: Node[] = [];

    const newDiamonds: string[] = [];
    const matchedDiamonds: string[] = [];

    for (const update of updates) {
        if (update.to.rank !== 1) {
            continue;
        }

        if (update.from.rank !== 1) {
            if (!update.to.joint) {
                newDiamonds.push(update.scoring);
            }
            else if (update.oldBestCurrentGolferCount === 1) {
                matchedDiamonds.push(update.scoring);
            }
        }
        else if (update.from.joint && !update.to.joint) {
            // Transition from golf to a new diamond.
            newDiamonds.push(update.scoring);
        }
    }

    if (newDiamonds.length) {
        popups.push(<div>
            <h3>Diamond Earned</h3>
            <p>New {newDiamonds.join('/')} 💎!</p>
        </div>);
    }

    if (matchedDiamonds.length) {
        popups.push(<div>
            <h3>Diamond Matched</h3>
            <p>Matched {matchedDiamonds.join('/')} 💎!</p>
        </div>);
    }

    return popups;
};

let lastSubmittedCode = '';
export function getLastSubmittedCode(){
    return lastSubmittedCode;
}

export async function submit(
    editor: any,
    // eslint-disable-next-line no-unused-vars
    updateReadonlyPanels: (d: ReadonlyPanelsData) => void,
): Promise<boolean> {
    if (!editor) return false;
    $('h2').innerText = '…';
    $('#status').className = 'grey';
    $$('canvas').forEach(e => e.remove());

    const code = editor.state.doc.toString();
    lastSubmittedCode = code;
    const codeLang = lang;
    const submissionID = ++latestSubmissionID;

    const res  = await fetch('/solution', {
        method: 'POST',
        body: JSON.stringify({ code, hole, lang }),
    });

    if (res.status != 200) {
        alert('Error ' + res.status);
        return false;
    }

    const data = await res.json() as SubmitResponse;
    savedInDB = data.logged_in;

    if (submissionID != latestSubmissionID)
        return false;

    const pass = data.runs.every(r => r.pass);
    $('main')?.classList.remove('pass');
    $('main')?.classList.remove('fail');
    $('main')?.classList.add(pass ? 'pass' : 'fail');
    $('main')?.classList.add('lastSubmittedCode');
    if (pass) {
        for (const i of [0, 1] as const) {
            const solutionCode = getSolutionCode(codeLang, i);
            if (!solutionCode || getScoring(code, i) <= getScoring(solutionCode, i)) {
                solutions[i][codeLang] = code;

                // Don't need to keep solution in local storage because it's
                // stored on the site. This prevents conflicts when the
                // solution is improved on another browser.
                if (savedInDB && localStorage.getItem(getAutoSaveKey(codeLang, i)) == code)
                    localStorage.removeItem(getAutoSaveKey(codeLang, i));
            }
        }
    }

    for (const i of [0, 1] as const) {
        const key = getAutoSaveKey(codeLang, i);
        if (savedInDB) {
            // If the auto-saved code matches either solution, remove it to
            // avoid prompting the user to restore it.
            const autoSaveCode = localStorage.getItem(key);
            for (const j of [0, 1] as const) {
                if (getSolutionCode(codeLang, j) == autoSaveCode)
                    localStorage.removeItem(key);
            }
        }
        else if (getSolutionCode(codeLang, i)) {
            // Autosave the best solution for each scoring metric, but don't
            // save two copies of the same solution, because that can lead to
            // the solution picker being show unnecessarily.
            if ((i == 0 || getSolutionCode(codeLang, 0) != getSolutionCode(codeLang, i)) &&
                getSolutionCode(codeLang, i) !== langs[codeLang].example)
                localStorage.setItem(key, getSolutionCode(codeLang, i));
            else
                localStorage.removeItem(key);
        }
    }

    // Automatically switch to the solution whose code matches the current
    // code after a new solution is submitted. Don't change scoring,
    // refreshScores will update the solution picker.
    if (pass && getSolutionCode(codeLang, solution) != code &&
        getSolutionCode(codeLang, getOtherScoring(solution)) == code)
        setSolution(getOtherScoring(solution));

    // Update the restore link visibility, after possibly changing the active
    // solution.
    updateRestoreLinkVisibility(editor);

    function showRun(run: Run) {
        updateReadonlyPanels({
            Pass: run.pass,
            Argv: run.args,
            Exp: run.answer,
            Err: run.stderr,
            Out: run.stdout,
            OutputDelimiter: run.output_delimiter,
            MultisetItemDelimiter: run.multiset_item_delimiter,
        });

        const ms = Math.round(run.time_ns / 10**6);
        // Only show runtime if it's more than 1000ms, for a bit less clutter in general.
        $('#runtime').innerText = ms > 1000 ? `(${ms}ms)` : '';

        // 3rd party integrations.
        let thirdParty = '';
        if (lang == 'hexagony') {
            const payload = LZString.compressToBase64(JSON.stringify({
                code, input: run.args.join('\0') + '\0', inputMode: 'raw' }));

            thirdParty = <a href={'//hexagony.net#lz' + payload}>
                Run on Hexagony.net
            </a>;
        }
        $('#thirdParty').replaceChildren(thirdParty);
    }

    // Default run: first failing non-timeout, else first timeout, else last overall.
    let defaultRunIndex = data.runs.findIndex(run => !run.pass && !run.timeout);
    if (defaultRunIndex === -1)
        defaultRunIndex = data.runs.findIndex(run => !run.pass);
    if (defaultRunIndex === -1)
        defaultRunIndex = data.runs.length - 1;

    const btns = data.runs.map((run, i) => {
        const [emoji, label] = run.pass ? ['😀', 'Pass']
            : run.timeout ? ['⏱️', 'Timeout']
                : ['☹️', 'Fail'];
        const btn = (
            <button class="run-result-btn" aria-label={`Run ${i + 1}: ${label}`}>
                {emoji}
            </button>
        );
        function onPickRun() {
            showRun(run);
            $$<HTMLButtonElement>('.run-result-btn').forEach(b => b.disabled = false);
            $('#pass-fail-msg').innerText = label;
            btn.disabled = true;
        };
        btn.addEventListener('click', onPickRun);
        return btn;
    });

    $('h2').replaceChildren(
        <span id="run-result-btns">{btns}</span>,
        <span id="pass-fail-msg"></span>,
        <span id="runtime"></span>,
    );

    btns[defaultRunIndex].click();

    $('#status').className = pass ? 'green' : 'red';

    // Show popups.
    $('#popups').replaceChildren(
        ...scorePopups(data.rank_updates),
        ...diamondPopups(data.rank_updates),
        ...data.cheevos.map(c => <div>
            <h3>Achievement Earned!</h3>
            { c.emoji }<p>{ c.name }</p>
        </div>));

    refreshScores(editor);

    return pass;
}

export function updateLocalStorage(code: string) {
    // Avoid future conflicts by only storing code locally that's
    // different from the server's copy.
    const serverCode = getCurrentSolutionCode();
    const key = getAutoSaveKey(getLang(), getSolution());
    const hadLocalStorage = localStorage.getItem(key) !== null;
    const wantLocalStorage = code && (code !== serverCode || !getSavedInDB()) && code !== langs[getLang()].example;

    if (wantLocalStorage)
        localStorage.setItem(key, code);
    else
        localStorage.removeItem(key);

    if (wantLocalStorage !== hadLocalStorage && serverCode === '') {
        // We may be adding or removing a language slot.
        updateLangPicker();
    }
}

export function updateRestoreLinkVisibility(editor: any) {
    const restoreLink = $('#restoreLink');
    if (restoreLink instanceof HTMLAnchorElement) {
        const serverCode = getSolutionCode(lang, solution);
        const sampleCode = langs[lang].example;
        const currentCode = editor?.state.doc.toString();
        restoreLink.classList.toggle('hide',
            (!serverCode && currentCode !== sampleCode) || currentCode === serverCode);
        restoreLink.textContent =
            currentCode === sampleCode ? 'Clear sample code' : 'Restore solution';
    }
}

export function setCodeForLangAndSolution(editor: any) {
    if (solution != 0 && getSolutionCode(lang, 0) == getSolutionCode(lang, 1)) {
        const autoSave0 = localStorage.getItem(getAutoSaveKey(lang, 0));
        const autoSave1 = localStorage.getItem(getAutoSaveKey(lang, 1));
        if (autoSave0 && !autoSave1)
            setSolution(0);
    }

    setState(localStorage.getItem(getAutoSaveKey(lang, solution)) ||
        getSolutionCode(lang, solution) || langs[lang].example, editor);

    if (lang == 'assembly') scoring = 0;
    const charsTab = $('#scoringTabs a:last-child');
    if (charsTab)
        charsTab.classList.toggle('hide', lang == 'assembly');

    refreshScores(editor);

    $$('#info-container .info').forEach(
        i => i.classList.toggle('hide', !i.classList.contains(lang)));
}

export async function populateScores(editor: any) {
    // Populate the rankings table.
    if (!$('#scores')) return;
    const scoringID = scorings[scoring].toLowerCase();
    const path      = `/${hole}/${lang}/${scoringID}`;
    const view      = $('#rankingsView a:not([href])').innerText.trim().toLowerCase();
    const res       = await fetch(`/api/mini-rankings${path}/${view}` + (tabLayout ? '?long=1' : ''));
    const rows      = res.ok ? await res.json() : [];
    const colspan   = lang == 'assembly' ? 3 : 4;

    $<HTMLAnchorElement>('#allLink').href = '/rankings/holes' + path;

    $('#scores').replaceChildren(<tbody class={scoringID}>{
        // Rows.
        rows.length || !tabLayout ? rows.map((r: any) => <tr class={r.me ? 'me' : ''}>
            <td>{r.rank}<sup>{ord(r.rank)}</sup></td>
            <td>
                <a href={`/golfers/${r.golfer.name}`}>
                    <img src={`/golfers/${r.golfer.name}/avatar/48`}/>
                    <span>{r.golfer.name}</span>
                </a>
            </td>
            <td title={tooltip(r, 'Bytes')}>
                {scoringID != 'bytes' ? comma(r.bytes) :
                    <a href={`/golfers/${r.golfer.name}/${hole}/${lang}/bytes`}>
                        <span>{comma(r.bytes)}</span>
                    </a>}
            </td>
            {lang == 'assembly' ? '' : <td title={tooltip(r, 'Chars')}>
                {scoringID != 'chars' ? comma(r.chars) :
                    <a href={`/golfers/${r.golfer.name}/${hole}/${lang}/chars`}>
                        <span>{comma(r.chars)}</span>
                    </a>}
            </td>}
        </tr>): <tr><td colspan={colspan}>(Empty)</td></tr>
    }</tbody>);

    // Scroll the rankings to the top or the "me" row if applicable.
    const me            = $('.me');
    const scoresWrapper = $('#scores-wrapper');
    scoresWrapper.scrollTop = (view === 'me' && me)
        ? me.offsetTop + (me.offsetHeight / 2) - (scoresWrapper.offsetHeight / 2)
        : 0;

    $$<HTMLAnchorElement>('#scoringTabs a').forEach((tab, i) => {
        if (tab.innerText == scorings[scoring]) {
            tab.removeAttribute('href');
            tab.onclick = () => {};
        }
        else {
            tab.href = '';
            tab.onclick = e  => {
                e.preventDefault();
                scoring = i as 0 | 1;
                localStorage.setItem('scoring', scorings[scoring]);
                refreshScores(editor);
            };
        }
    });
}

export function getScorings(tr: any, editor: any) {
    const code = tr.state.doc.toString();
    const total: {byte?: number, char?: number} = {};
    const selection: {byte?: number, char?: number} = {};

    if (getLang() == 'assembly')
        total.byte = (editor.state.field(ASMStateField) as any).head.length();
    else {
        total.byte = byteLen(code);
        total.char = charLen(code);

        if (tr.selection) {
            selection.byte = 0;
            selection.char = 0;

            for (const range of tr.selection.ranges) {
                const contents = code.slice(range.from, range.to);
                selection.byte += byteLen(contents);
                selection.char += charLen(contents);
            }
        }
    }

    return (selection.byte || selection.char) ? {total, selection} : {total};
}

export function ctrlEnter(func: Function) {
    return function (e: KeyboardEvent) {
        if ((e.ctrlKey || e.metaKey) && e.key == 'Enter') {
            return func();
        }
    };
}
