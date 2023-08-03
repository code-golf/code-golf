import { ASMStateField }                       from '@defasm/codemirror';
import { $, $$, byteLen, charLen, comma, ord } from './_util';
import { Vim }                                 from '@replit/codemirror-vim';
import { EditorState, EditorView, extensions } from './_codemirror';
import pbm                                     from './_pbm';
import LZString                                from 'lz-string';

let tabLayout: boolean = false;

export function init(_tabLayout: boolean, setSolution: any, setCodeForLangAndSolution: any, updateReadonlyPanels: any, getEditor: () => any) {
    tabLayout = _tabLayout;
    const closuredSubmit = () => submit(getEditor(), updateReadonlyPanels);
    window.onkeydown = e => (e.ctrlKey || e.metaKey) && e.key == 'Enter'
        ? closuredSubmit()
        : undefined;
    if (vimMode) Vim.defineEx('write', 'w', closuredSubmit);

    (onhashchange = () => {
        const hashLang = location.hash.slice(1) || localStorage.getItem('lang');

        // Kick 'em to Python if we don't know the chosen language, or if there is no given language.
        lang = hashLang && langs[hashLang] ? hashLang : 'python';

        // Assembly only has bytes.
        if (lang == 'assembly')
            setSolution(0);

        localStorage.setItem('lang', lang);

        history.replaceState(null, '', '#' + lang);

        const editor = getEditor();
        if (tabLayout) {
            $('#hole-lang summary').innerText = langs[lang].name;
            refreshScores(editor);
        }
        setCodeForLangAndSolution(editor);
    })();

    $('dialog [name=text]').addEventListener('input', (e: Event) => {
        const target = e.target as HTMLInputElement;
        target.form!.confirm.toggleAttribute('disabled',
            target.value !== target.placeholder);
    });
}

export function initDeleteBtn(deleteBtn: HTMLElement | undefined, langs: any) {
    deleteBtn?.addEventListener('click', () => {
        $('#delete-dialog b').innerText = langs[lang].name;
        $<HTMLInputElement>('#delete-dialog [name=lang]').value = lang;
        $<HTMLInputElement>('#delete-dialog [name=text]').value = '';
        $<HTMLDialogElement>('#delete-dialog').showModal();
    });
}

export function initCopyJSONBtn(copyBtn: HTMLElement | undefined) {
    copyBtn?.addEventListener('click', () =>
        navigator.clipboard.writeText($('#data').innerText));
}

export const langs = JSON.parse($('#langs').innerText);
const sortedLangs  =
    Object.values(langs).sort((a: any, b: any) => a.name.localeCompare(b.name));
let lang: string = '';

export function getLang() {
    return lang;
}

const experimental = JSON.parse($('#experimental').innerText);
export const hole         = decodeURI(location.pathname.slice(1));
const scorings     = ['Bytes', 'Chars'];
const solutions    = JSON.parse($('#solutions').innerText);

const vimMode = JSON.parse($('#keymap').innerText) === 'vim';
const vimModeExtensions = vimMode ? [extensions.vim] : [];

const baseExtensions = [...vimModeExtensions, ...extensions.base];

let latestSubmissionID = 0;
let solution = scorings.indexOf(localStorage.getItem('solution') ?? 'Bytes') as 0 | 1;
let scoring  = scorings.indexOf(localStorage.getItem('scoring')  ?? 'Bytes') as 0 | 1;

let hideDeleteBtn: boolean = false;

// The savedInDB state is used to avoid saving solutions in localStorage when
// those solutions match the solutions in the database. It's used to avoid
// restoring a solution from localStorage when the user has improved that
// solution on a different browser. Assume the user is logged-in by default
// for non-experimental holes. At this point, it doesn't matter whether the
// user is actually logged-in, because solutions dictionaries will be empty
// for users who aren't logged-in, so the savedInDB state won't be used.
// By the time they are non-empty, the savedInDB state will have been updated.
let savedInDB = !experimental;

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

function getSolutionCode(lang: string, solution: 0 | 1) {
    return lang in solutions[solution] ? solutions[solution][lang] : '';
}

/**
 * Get the code corresponding to the current lang and solution (bytes/chars)
 */
export function getCurrentSolutionCode() {
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
                    ? [] : extensions.bracketMatching,
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
    // Populate the language picker with accurate stroke counts.
    $('#picker').replaceChildren(...sortedLangs.map((l: any) => {
        const tab = <a href={l.id == lang ? null : '#'+l.id}>{l.name}</a>;

        if (getSolutionCode(l.id, 0)) {
            const bytes = byteLen(getSolutionCode(l.id, 0));
            const chars = charLen(getSolutionCode(l.id, 1));

            let text = comma(bytes);
            if (chars && bytes != chars) text += '/' + comma(chars);

            tab.append(' ', <sup>{text}</sup>);
        }

        return tab;
    }));
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

    // Hide the delete button for exp holes or if we have no solutions.
    hideDeleteBtn = experimental || (!dbBytes && !dbChars);
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
    from: RankFromTo,
    to: RankFromTo,
    beat: number | null,
    oldBestJoint: boolean | null,
    oldBestStrokes: number | null,
}

export interface Run {
    answer: string,
    args: string[],
    exit_code: number,
    pass: boolean,
    stderr: string,
    stdout: string,
    time_ns: number,
    timeout: boolean
}

export interface ReadonlyPanelsData {
    Pass: boolean,
    Out: string,
    Exp: string,
    Err: string,
    Argv: string[],
}

export interface SubmitResponse {
    Pass: boolean,
    Cheevos: {
        emoji: string,
        name: string
    }[],
    LoggedIn: boolean,
    RankUpdates: RankUpdate[],
    runs: Run[]
}

const makeSingular = (strokes: number, units: string) =>
    strokes == 1 ? units.substring(0, units.length - 1) : units;

const scorePopups = (updates: RankUpdate[]) => {
    const popups: Node[] = [];

    const strokes = [0, 0];
    const points = [0, 0];

    for (const i of [0, 1] as const) {
        const update = updates[i];
        if (update.from.strokes && update.to.strokes) {
            strokes[i] = update.from.strokes - update.to.strokes;
            if (update.oldBestStrokes) {
                const newBest = Math.min(update.oldBestStrokes, update.to.strokes);
                points[i] = Math.round(newBest / update.to.strokes * 1000) - Math.round(update.oldBestStrokes / update.from.strokes * 1000);
            }
        }
    }

    if (strokes[0] > 0 || strokes[1] > 0) {
        let amount = '';
        if (strokes[0] > 0 && strokes[0] == strokes[1]) {
            const delta = strokes[0];
            let units = '';
            for (const i of [0, 1] as const) {
                units += (i == 1 ? '/' : '') + makeSingular(delta, updates[i].scoring);
            }

            amount = `${delta} ${units}`;
        }
        else {
            for (const i of [0, 1] as const) {
                if (strokes[i] > 0) {
                    amount += (i == 1 && strokes[0] > 0 ? '/' : '') + `${strokes[i]} ${makeSingular(strokes[i], updates[i].scoring)}`;
                }
            }
        }

        const pointsNodes: Node[] = [];
        if (points[0] > 0 && points[1] == points[0]) {
            pointsNodes.push(<p>Earned {points[0]} {makeSingular(points[0], 'points')} for {updates[0].scoring}/{updates[1].scoring}</p>);
        }
        else {
            for (const i of [0, 1] as const) {
                if (points[i] > 0) {
                    pointsNodes.push(<p>Earned {points[i]} {makeSingular(points[i], 'points')} for {updates[i].scoring}</p>);
                }
            }
        }

        popups.push(<div>
            <h3>Score Improved</h3>
            <p>Saved {amount}</p>
            {pointsNodes}
        </div>);
    }

    return popups;
};

const diamondPopups = (updates: RankUpdate[]) => {
    const popups: Node[] = [];

    const newDiamonds: string[] = [];
    const matchedDiamonds: string[] = [];

    for (const i of [0, 1] as const) {
        const update = updates[i];
        if (update.from.rank != 1 && update.to.rank == 1) {
            if (!update.to.joint) {
                newDiamonds.push(update.scoring);
            }
            else if (!update.oldBestJoint) {
                matchedDiamonds.push(update.scoring);
            }
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

export async function submit(
    editor: any,
    // eslint-disable-next-line no-unused-vars
    updateReadonlyPanels: (d: ReadonlyPanelsData) => void,
) {
    if (!editor) return;
    $('h2').innerText = '…';
    $('#status').className = 'grey';
    $$('canvas').forEach(e => e.remove());

    const code = editor.state.doc.toString();
    const codeLang = lang;
    const submissionID = ++latestSubmissionID;

    const res  = await fetch('/solution', {
        method: 'POST',
        body: JSON.stringify({
            Code: code,
            Hole: hole,
            Lang: lang,
        }),
    });

    if (res.status != 200) {
        alert('Error ' + res.status);
        return;
    }

    const data = await res.json() as SubmitResponse;
    savedInDB = data.LoggedIn && !experimental;

    if (submissionID != latestSubmissionID)
        return;

    if (data.Pass) {
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
    if (data.Pass && getSolutionCode(codeLang, solution) != code &&
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
        });

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

        if (hole == 'julia-set')
            $('main').append(pbm(run.answer) as Node, pbm(run.stdout) ?? [] as any);
    }

    // Default run: first failing, else last overall.
    const defaultRun = data.runs.find(run => !run.pass) ?? data.runs[data.runs.length-1];

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
            btn.disabled = true;
        };
        if (run === defaultRun) onPickRun();
        btn.addEventListener('click', onPickRun);
        return btn;
    });

    $('h2').replaceWith(<h2>
        {data.Pass ? 'Pass' : 'Fail'}
        <span class="btns">{btns}</span>
    </h2>);

    showRun(data.runs[0]);

    $('#status').className = data.Pass ? 'green' : 'red';

    // Show cheevos.
    $('#popups').replaceChildren(...scorePopups(data.RankUpdates),
        ...diamondPopups(data.RankUpdates),
        ...data.Cheevos.map(c => <div>
            <h3>Achievement Earned!</h3>
            { c.emoji }<p>{ c.name }</p>
        </div>));

    refreshScores(editor);
}

export function updateRestoreLinkVisibility(editor: any) {
    const serverCode = getSolutionCode(lang, solution);
    $('#restoreLink')?.classList.toggle('hide',
        !serverCode || editor?.state.doc.toString() == serverCode);
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

    $$('main .info').forEach(
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

    $<HTMLAnchorElement>('#allLink').href = '/rankings/holes' + path;

    $('#scores').replaceChildren(<tbody class={scoringID}>{
        // Rows.
        rows.length || !tabLayout ? rows.map((r: any) => <tr class={r.me ? 'me' : ''}>
            <td>{r.rank}<sup>{ord(r.rank)}</sup></td>
            <td>
                <a href={`/golfers/${r.golfer.name}`}>
                    <img src={`//avatars.githubusercontent.com/${r.golfer.name}?s=48`}/>
                    <span>{r.golfer.name}</span>
                </a>
            </td>
            <td data-tooltip={tooltip(r, 'Bytes')}>{comma(r.bytes)}</td>
            <td data-tooltip={tooltip(r, 'Chars')}>{comma(r.chars)}</td>
        </tr>): <tr><td colspan="4">(Empty)</td></tr>
    }{
        // Padding.
        tabLayout ? [] : [...Array(7 - rows.length).keys()].map(() =>
            <tr><td colspan="4">&nbsp;</td></tr>)
    }</tbody>);

    if (tabLayout) {
        if (view === 'me')
            $('.me')?.scrollIntoView({block: 'center'});
        else
            $('#scores-wrapper').scrollTop = 0;
    }

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
