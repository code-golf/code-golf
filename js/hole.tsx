import { EditorView }   from './_codemirror';
import diffTable        from './_diff';
import { $, $$, comma } from './_util';
import {
    init, langs, getLang, hole, getAutoSaveKey, setSolution, getSolution,
    setCode, refreshScores, submit, getSavedInDB, updateRestoreLinkVisibility,
    ReadonlyPanelsData, setCodeForLangAndSolution, getCurrentSolutionCode,
    initDeleteBtn, initCopyJSONBtn, getScorings,
} from './_hole-common';

const editor = new EditorView({
    dispatch: tr => {
        const result = editor.update([tr]) as unknown;

        const code = tr.state.doc.toString();
        const scorings: {total: {byte?: number, char?: number}, selection?: {byte?: number, char?: number}} = getScorings(tr, editor);
        const scoringKeys = ['byte', 'char'] as const;

        function formatScore(scoring: any) {
            return scoringKeys
                .filter(s => s in scoring)
                .map(s => `${comma(scoring[s])} ${s}${scoring[s] != 1 ? 's' : ''}`)
                .join(', ');
        }

        $('#strokes').innerText = scorings.selection
            ? `${formatScore(scorings.total)} (${formatScore(scorings.selection)} selected)`
            : formatScore(scorings.total);

        // Avoid future conflicts by only storing code locally that's
        // different from the server's copy.
        const serverCode = getCurrentSolutionCode();

        const key = getAutoSaveKey(getLang(), getSolution());
        if (code && (code !== serverCode || !getSavedInDB()) && code !== langs[getLang()].example)
            localStorage.setItem(key, code);
        else
            localStorage.removeItem(key);

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
$('#runBtn').onclick = () => submit(editor, updateReadonlyPanels);

initCopyJSONBtn($('#copy'));
initDeleteBtn($('#deleteBtn'), langs);

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
    $('#out div').innerText = data.Out;

    const diff = diffTable(hole, data.Exp, data.Out, data.Argv);
    $('#diff-content').replaceChildren(diff);
    $('#diff').classList.toggle('hide', !diff);
}
