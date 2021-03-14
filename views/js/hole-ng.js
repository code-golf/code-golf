import { EditorView, EditorState, extensions, languages } from '/editor.js';

const langs  = JSON.parse(document.querySelector('#langs').innerText);
const select = document.querySelector('select');
const editor = new EditorView({ parent: document.querySelector('#editor') });
editor.contentDOM.setAttribute("data-gramm", "false"); // Disable Grammarly

// Switch Lang
const switchLang = onhashchange = () => {
    let lang = location.hash.slice(1);

    // Kick 'em to Python if we don't know the chosen language.
    if (!langs[lang]) {
        lang = 'python';

        history.replaceState(null, '', '#' + lang);
    }

    select.value = lang;

    editor.setState(
        EditorState.create({
            doc:        langs[lang].example,
            extensions: [...extensions, languages[lang] || []],
        }),
    );
};

select.onchange = e => {
    history.replaceState(null, '', '#' + e.target.value);
    switchLang();
}

switchLang();

// Run Code
const runCode = document.querySelector('#run a').onclick = () => {
    alert(editor.state.doc);
};

onkeydown = e => e.ctrlKey && e.key == 'Enter' ? runCode() : undefined;
