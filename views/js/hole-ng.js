// TODO Set #rankings height to editor on expand/shrink.

import { darkThemeExtensions, defaultExtensions, EditorView, EditorState, languages, LZString } from '/editor.js';

let lang;
let result = {};
let scoring = 'bytes';

const all         = document.querySelector('#all');
const hole        = decodeURI(location.pathname.slice(4));
const langs       = JSON.parse(document.querySelector('#langs').innerText);
const rankings    = document.querySelector('#rankings');
const scoringTabs = document.querySelectorAll('#scoringTabs a');
const select      = document.querySelector('select');
const solutions   = JSON.parse(document.querySelector('#solutions').innerText);
const status      = document.querySelector('#status');
const statusView  = new EditorView({ parent: status });
const statusTabs  = document.querySelectorAll('#statusTabs a');
const strokes     = document.querySelector('#strokes');
const thirdParty  = document.querySelector('#thirdParty');
const editor      = new EditorView({
    dispatch: (tr) => {
        const result = editor.update([tr]);
        let scorings = {};

        if (lang == 'assembly')
            scorings.byte = editor['asm-bytes'];
        else {
            const code = [...tr.state.doc].join('');

            scorings.byte = new TextEncoder().encode(code).length;
            scorings.char = strlen(code);
        }

        strokes.innerText = Object.keys(scorings).map(
            s => `${scorings[s]} ${s}${scorings[s] != 1 ? 's' : ''}`
        ).join(', ');

        return result;
    },
    parent: document.querySelector('#editor'),
});

editor.contentDOM.setAttribute("data-gramm", "false"); // Disable Grammarly

const darkMode = document.documentElement.classList.contains('dark');
const extensions = darkMode ?
    [...darkThemeExtensions, ...defaultExtensions] :
    defaultExtensions;

// Update UI.
async function update() {
    // From left to right... update lang select.
    const svg = document.querySelector('#' + lang);
    svg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    const color = (darkMode ? '#1e2124' : '#fdfdfd').replaceAll('#', '%23');
    const data = svg.outerHTML.replaceAll('currentColor', color)
        .replaceAll('#', '%23').replaceAll('"', "'");

    select.style.background = `url("data:image/svg+xml,${data}") no-repeat left .5rem center/1rem auto, url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 4 5'><path d='M2 0L0 2h4zm0 5L0 3h4' fill='${color}'/></svg>") no-repeat right .5rem center/auto calc(100% - 1.25rem), var(--color)`;

    // Update scoring tabs.
    for (const tab of scoringTabs)
        if (tab.id == scoring)
            tab.removeAttribute('href');
        else
            tab.href = '';

    // Update "All" link.
    all.href = `/rankings/holes/${hole}/${lang}/${scoring}`;

    const res  = await fetch(`/scores/${hole}/${lang}/${scoring}/mini?ng=1`);
    const rows = res.ok ? await res.json() : [];

    let html = '';
    for (const r of rows)
        html += `
            <tr ${r.me ? 'class=me' : ''}>
                <td>${r.rank}<sup>${ord(r.rank)}</sup>
                <td>
                    <a href=/golfers/${r.login}>
                        <img src="//avatars.githubusercontent.com/${r.login}?s=24">
                        <span>${r.login}</span>
                    </a>
                <td class=right>${r[scoring].toLocaleString('en')}`;

    rankings.innerHTML =
        html + '<tr><td colspan=3>&nbsp;'.repeat(15 - rows.length);
}

// Switch scoring
for (const tab of scoringTabs)
    tab.onclick = e => { e.preventDefault(); scoring = tab.id; update() };

// Switch status
for (const tab of statusTabs)
    tab.onclick = e => {
        e.preventDefault();

        for (const tab of statusTabs)
            tab.href = '';

        tab.removeAttribute('href');

        let ext = extensions;
        if (tab.id == 'diff')
            ext.push(languages.diff);

        statusView.setState(EditorState.create({
            doc: result[tab.id], extensions: extensions }));
    };

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
            doc:        solutions[lang]?.[scoring] ?? langs[lang].example,
            extensions: [...extensions, languages[lang] || []],
        }),
    );

    if (lang == 'assembly') {
        scoring = 'bytes';
        scoringTabs[1].style.display = 'none';
    }
    else
        scoringTabs[1].style.display = '';

    for (const info of document.querySelectorAll('main .info'))
        info.style.display = info.classList.contains(lang) ? 'block' : '';

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
const runCode = document.querySelector('#run a').onclick = async () => {
    status.style.display = 'none';

    const code = [...editor.state.doc].join('');
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

    const data = await res.json();

    status.style.background = data.Pass ? 'var(--green)' : 'var(--red)';

    document.querySelector('h2').innerText
        = data.Pass ? 'Pass 😀' : 'Fail ☹️';

    result = {
        arguments: data.Argv.join('\n'),
        diff:      data.Diff,
        errors:    data.Err,
        expected:  data.Exp,
        output:    data.Out,
    };

    // 3rd party integrations.
    if (lang == 'hexagony') {
        const url = '//hexagony.net#lz' + LZString.compressToBase64(
            JSON.stringify({ code, input: data.Argv.join('\0'), inputMode: 'raw' }));

        thirdParty.innerHTML =
            `<a href="${url}" target=_blank>Run on Hexagony.net</a>`;
    } else
        thirdParty.innerHTML = '';

    // Default to the "Diff" tab.
    statusTabs[1].click();

    status.style.display = 'grid';
};

onkeydown = e => (e.ctrlKey || e.metaKey) && e.key == 'Enter' ? runCode() : undefined;

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
