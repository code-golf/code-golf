import { EditorView }    from '../_codemirror';
import diffView          from '../_diff';
import { $, $$, amount } from '../_util';
import {
    init, hole, setSolution,
    setCode, refreshScores, submit, updateRestoreLinkVisibility,
    ReadonlyPanelsData, setCodeForLangAndSolution, getCurrentSolutionCode,
    initDeleteBtn, initCopyButtons, getScorings, describeSelectedChar,
    updateLocalStorage,
    ctrlEnter,
    getLastSubmittedCode,
    Scorings,
    getLang, setState, getHoleLangNotesContent, setHoleLangNotesCache,
} from '../_hole-common';
import UnprintableElement from '../_unprintable';

const editor = new EditorView({
    dispatch: tr => {
        const result = editor.update([tr]) as unknown;

        const code = tr.state.doc.toString();
        const scorings: {total: Scorings, selection?: Scorings} = getScorings(tr, editor);
        const scoringKeys = ['byte', 'char', 'stroke'] as const;

        $('main')?.classList.toggle('lastSubmittedCode', code === getLastSubmittedCode());

        function formatScore(scoring: any) {
            return scoringKeys
                .filter(s => s in scoring)
                .map(s => amount(scoring[s], s))
                .join(', ');
        }

        $('#strokes').innerText = (scorings.selection
            ? `${formatScore(scorings.total)} (${formatScore(scorings.selection)} selected)`
            : formatScore(scorings.total)) + describeSelectedChar(tr.state, scorings);

        updateLocalStorage(code);
        updateRestoreLinkVisibility(editor);

        return result;
    },
    parent: $('#editor'),
});

editor.contentDOM.setAttribute('data-gramm', 'false');  // Disable Grammarly.

init(false, setSolution, setCodeForLangAndSolution, updateReadonlyPanels, () => editor);

// Set/clear the hide-details cookie on details toggling.
$('#details').ontoggle = (e: Event) => document.cookie =
    'hide-details=;SameSite=Lax;Secure' + ((e.target as HTMLDetailsElement).open ? ';Max-Age=0' : '');

$('#restoreLink').onclick = e => {
    setCode(getCurrentSolutionCode(), editor);
    e.preventDefault();
};

// Wire submit to clicking a button and a keyboard shortcut.
const closuredSubmit = () => submit(editor, updateReadonlyPanels);
$('#runBtn').onclick = closuredSubmit;

const saveNotes = initNotes();
window.onkeydown = ctrlEnter(() => isNotesMode() ? saveNotes?.() : closuredSubmit());

initCopyButtons($$('[data-copy]'));
initDeleteBtn($('#deleteBtn'));

$$('#rankingsView a').forEach(a => a.onclick = e => {
    e.preventDefault();

    $$<HTMLAnchorElement>('#rankingsView a').forEach(a => a.href = '');
    a.removeAttribute('href');

    document.cookie =
        `rankings-view=${a.innerText.toLowerCase()};SameSite=Lax;Secure`;

    refreshScores(editor);
});

function updateReadonlyPanels(data: ReadonlyPanelsData) {
    // Hide arguments unless we have some.
    $('#arg div').replaceChildren(...data.Argv.map(a => <span>{a}</span>));
    $('#arg').classList.toggle('hide', !data.Argv.length);

    // Hide stderr if we're passing or have no stderr output.
    $('#err div').innerHTML = data.Err.replace(/\n/g, '<br>');
    $('#err').classList.toggle('hide', data.Pass || !data.Err);

    // Always show exp & out.
    $('#exp div').innerText = data.Exp;
    $('#out div').replaceChildren(UnprintableElement.escape(data.Out));

    const ignoreCase = JSON.parse($('#case-fold').innerText);
    const diff = diffView(hole, data.Exp, data.Out, data.Argv, ignoreCase, data.OutputDelimiter, data.MultisetItemDelimiter);
    $('#diff-content').replaceChildren(diff);
    $('#diff').classList.toggle('hide', !diff);
}

const isNotesMode = () => $('main').classList.contains('notes-mode');

// Notes per hole per lang, the same notes as in the tabs layout. The editor
// tabs switch between the code editor and the notes editor.
function initNotes() {
    const editorTabs = $('#editorTabs');
    if (!editorTabs) return;

    const [codeTab, notesTab] = $$<HTMLAnchorElement>('#editorTabs a');

    const isSponsor = JSON.parse($('#sponsor').innerText);
    const langs     = JSON.parse($('#languages').innerText);
    const saveBtn   = $<HTMLButtonElement>('#saveNotesBtn');
    let savedContent = '';

    // Only sponsors can create or update notes, the rest can still read & delete.
    if (!isSponsor) $('#notesStatus').innerText = 'Only sponsors can edit notes.';

    const notesEditor: EditorView = new EditorView({
        dispatch: tr => {
            const result = notesEditor.update([tr]) as unknown;
            const content = tr.state.doc.toString();
            saveBtn.disabled = content === savedContent || (!!content && !isSponsor);
            return result;
        },
        parent: $('#notes'),
    });

    notesEditor.contentDOM.setAttribute('data-gramm', 'false');  // Disable Grammarly.

    const refreshNotes = async () => {
        const lang = getLang();
        $('#notesTitle').innerText = `${langs[lang]?.name ?? lang} notes`;

        const content = await getHoleLangNotesContent(lang);

        // Ignore the response if the lang changed in the meantime.
        if (lang !== getLang()) return;

        savedContent = content;
        setState(content, notesEditor);
    };

    const setNotesMode = (on: boolean) => {
        if (on === isNotesMode()) return;

        $('main').classList.toggle('notes-mode', on);

        // The tabs sit above the top right corner of whichever editor is shown.
        $(on ? '#notes .header-actions' : '#editor .header-actions').append(editorTabs);

        // Like the other tabs, the selected one has no href.
        (on ? codeTab : notesTab).href = '';
        (on ? notesTab : codeTab).removeAttribute('href');

        (on ? notesEditor : editor).focus();
    };

    codeTab.onclick = e => {
        e.preventDefault();
        setNotesMode(false);
    };

    notesTab.onclick = e => {
        e.preventDefault();
        setNotesMode(true);
    };

    // Picking any lang, even the current one, goes back to the code.
    $('#picker').addEventListener('click', e => {
        if ((e.target as Element).closest('a')) setNotesMode(false);
    });
    $('#picker').addEventListener('change', () => setNotesMode(false));
    addEventListener('hashchange', () => setNotesMode(false));

    const save = async () => {
        if (saveBtn.disabled) return;
        saveBtn.disabled = true;

        const lang = getLang();
        const content = notesEditor.state.doc.toString();
        const resp = await fetch(
            `/api/notes/${hole}/${lang}`,
            content ? { method: 'PUT', body: content } : { method: 'DELETE' },
        );

        if (resp.status !== 204) {
            saveBtn.disabled = false;
            return;
        }

        setHoleLangNotesCache(lang, content);
        if (lang === getLang()) savedContent = content;
    };

    saveBtn.onclick = save;
    addEventListener('hashchange', refreshNotes);
    refreshNotes();

    return save;
}
