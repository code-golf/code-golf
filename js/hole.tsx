import { ASMStateField }                  from '@defasm/codemirror';
import { EditorView }                     from './_codemirror';
import diffTable                          from './_diff';
import { $, $$, byteLen, charLen, comma } from './_util';
import {
    init, langs, getLang, hole, getAutoSaveKey, setSolution, getSolution,
    setCode, refreshScores, submit, getSavedInDB, updateRestoreLinkVisibility,
    SubmitResponse, setCodeForLangAndSolution, getCurrentSolutionCode,
    initDeleteBtn, initCopyJSONBtn,
} from './_hole-common';

const editor = new EditorView({
    dispatch: tr => {
        const result = editor.update([tr]) as unknown;

        const code = tr.state.doc.toString();
        const scorings: {byte?: number, char?: number} = {};
        const scoringKeys = ['byte', 'char'] as const;

        if (getLang() == 'assembly')
            scorings.byte = (editor.state.field(ASMStateField) as any).head.length();
        else {
            scorings.byte = byteLen(code);
            scorings.char = charLen(code);
        }

        $('#strokes').innerText = scoringKeys
            .filter(s => s in scorings)
            .map(s => `${comma(scorings[s])} ${s}${scorings[s] != 1 ? 's' : ''}`)
            .join(', ');

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

    refreshScores(setCodeForLangAndSolution);
});

function updateReadonlyPanels(data: SubmitResponse) {
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
