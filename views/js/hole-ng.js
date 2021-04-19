// TODO Set #rankings height to editor on expand/shrink.

import { EditorView, EditorState, extensions, languages } from '/editor.js';

let lang;
let scoring = 'bytes';

const all       = document.querySelector('#all');
const hole      = decodeURI(location.pathname.slice(4));
const langs     = JSON.parse(document.querySelector('#langs').innerText);
const rankings  = document.querySelector('#rankings table');
const select    = document.querySelector('select');
const solutions = JSON.parse(document.querySelector('#solutions').innerText);
const strokes   = document.querySelector('#strokes');
const tabs      = document.querySelectorAll('.tabs a');
const editor    = new EditorView({
    dispatch: (tr) => {
        const code  = [...tr.state.doc].join('');
        const bytes = new TextEncoder().encode(code).length;
        const chars = strlen(code);

        strokes.innerText =
            `${bytes} byte${bytes != 1 ? 's' : ''}, ` +
            `${chars} char${chars != 1 ? 's' : ''}`;

        return editor.update([tr]);
    },
    parent: document.querySelector('#editor'),
});

editor.contentDOM.setAttribute("data-gramm", "false"); // Disable Grammarly

// Update UI.
async function update() {
    // From left to right... update lang select.
    const svg = document.querySelector('#' + lang);
    svg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    const data = svg.outerHTML.replaceAll('currentColor', '#fdfdfd')
        .replaceAll('#', '%23').replaceAll('"', "'");

    select.style.background = `url("data:image/svg+xml,${data}") no-repeat left .5rem center/1rem auto, url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 4 5'><path d='M2 0L0 2h4zm0 5L0 3h4' fill='%23fff'/></svg>") no-repeat right .5rem center/auto calc(100% - 1.25rem), var(--color)`;

    // Update scoring tabs.
    for (const tab of tabs)
        if (tab.id == scoring)
            tab.removeAttribute('href');
        else
            tab.href = '';

    // Update "All" link.
    all.href = `/rankings/holes/${hole}/${lang}/${scoring}`;

    const res = await fetch(`/scores/${hole}/${lang}/${scoring}/mini`);

    let html = '';
    for (const r of res.ok ? await res.json() : [])
        html += `
            <tr ${r.me ? 'class=me' : ''}>
                <td>${r.rank}<sup>${ord(r.rank)}</sup>
                <td>
                    <a href=/golfers/${r.login}>
                        <img src="//avatars.githubusercontent.com/${r.login}?s=24">
                        <span>${r.login}</span>
                    </a>
                <td class=right>${r[scoring].toLocaleString('en')}`;

    rankings.innerHTML = html;
}

// Switch scoring
for (const tab of tabs)
    tab.onclick = e => { e.preventDefault(); scoring = tab.id; update() };

// Switch lang
const switchLang = onhashchange = () => {
    lang = location.hash.slice(1);

    // Kick 'em to Python if we don't know the chosen language.
    if (!langs[lang]) {
        lang = 'python';

        history.replaceState(null, '', '#' + lang);
    }

    select.value = lang;

    editor.setState(
        EditorState.create({
            doc:        solutions[lang].bytes ?? langs[lang].example,
            extensions: [...extensions, languages[lang] || []],
        }),
    );

    // Dispatch to update strokes.
    editor.dispatch();

    update();
};

select.onchange = () => {
    history.replaceState(null, '', '#' + select.value);
    switchLang();
}

switchLang();

// Run Code
const runCode = document.querySelector('#run a').onclick = () => {
    alert(editor.state.doc);
};

onkeydown = e => e.ctrlKey && e.key == 'Enter' ? runCode() : undefined;

// Adapted from https://codegolf.stackexchange.com/a/119563
const ord = i => [, 'st', 'nd', 'rd'][i % 100 >> 3 ^ 1 && i % 10] || 'th';

// Adapted from https://mths.be/punycode
function strlen(str) {
    let i = 0, len = 0;

    while (i < str.length) {
        const value = str.charCodeAt(i++);

        if (value >= 0xD800 && value <= 0xDBFF && i < str.length) {
            // It's a high surrogate, and there is a next character.
            const extra = str.charCodeAt(i++);

            // Low surrogate.
            if ((extra & 0xFC00) == 0xDC00) {
                len++;
            } else {
                // It's an unmatched surrogate; only append this code unit, in case the
                // next code unit is the high surrogate of a surrogate pair.
                len++;
                i--;
            }
        } else {
            len++;
        }
    }

    return len;
}
