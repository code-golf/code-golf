import { ASMStateField }                       from '@defasm/codemirror';
import LZString                                from 'lz-string';
import { EditorState, EditorView, extensions } from './_codemirror.js';
import pbm                                     from './_pbm.js';
import strlen                                  from './_strlen.js';
import diffHTML                                from './_diff';

const all         = document.querySelector('#all');
const hole        = decodeURI(location.pathname.slice(4));
const langs       = JSON.parse(document.querySelector('#langs').innerText);
const rankings    = document.querySelector('#rankings');
const scoringTabs = document.querySelectorAll('#scoringTabs a');
const select      = document.querySelector('select');
const solutions   = JSON.parse(document.querySelector('#solutions').innerText);
const status      = document.querySelector('#status');
const statusH2    = document.querySelector('#status h2');
const strokes     = document.querySelector('#strokes');
const thirdParty  = document.querySelector('#thirdParty');
const diffContent = document.querySelector('#diff-content');

const darkMode = matchMedia(JSON.parse(document.querySelector(
    '#darkModeMediaQuery').innerText)).matches;

const baseExtensions =
    darkMode ? [...extensions.dark, ...extensions.base] : extensions.base;

let lang, scoring = 'bytes';

const editor = new EditorView({
    dispatch: tr => {
        const result = editor.update([tr]);
        const scorings = {};

        if (lang == 'assembly')
            scorings.byte = editor.state.field(ASMStateField).head.length();
        else {
            const code = [...tr.state.doc].join('');

            scorings.byte = new TextEncoder().encode(code).length;
            scorings.char = strlen(code);
        }

        strokes.innerText = Object.keys(scorings)
            .map(s => `${scorings[s]} ${s}${scorings[s] != 1 ? 's' : ''}`)
            .join(', ');

        return result;
    },
    parent: document.querySelector('#editor'),
});

editor.contentDOM.setAttribute('data-gramm', 'false');  // Disable Grammarly.

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

    rankings.replaceChildren(<tbody>{
        // Rows.
        rows.map(r => <tr class={ r.me ? 'me' : '' }>
            <td>{ r.rank }<sup>{ ord(r.rank) }</sup></td>
            <td>
                <a href={ `/golfers/${r.login}` }>
                    <img src={ `//avatars.githubusercontent.com/${r.login}?s=24` }/>
                    <span>{ r.login }</span>
                </a>
            </td>
            <td class="right">{ r[scoring].toLocaleString('en') }</td>
        </tr>)
    }{
        // Padding.
        [...Array(15 - rows.length).keys()].map(() =>
            <tr><td colspan="3">&nbsp;</td></tr>)
    }</tbody>);
}

// Switch scoring
for (const tab of scoringTabs)
    tab.onclick = e => { e.preventDefault(); scoring = tab.id; update() };

// Switch lang
const switchLang = onhashchange = () => {
    lang = location.hash.slice(1) || localStorage.getItem('lang');

    // Kick 'em to Python if we don't know the chosen language.
    if (!langs[lang])
        lang = 'python';

    select.value = lang;
    localStorage.setItem('lang', lang);
    window.history.replaceState(null, '', '#' + lang);

    editor.setState(
        EditorState.create({
            doc:        solutions[lang]?.[scoring] ?? langs[lang].example,
            extensions: [
                ...baseExtensions,

                extensions[lang] || [],

                // These languages shouldn't match brackets.
                ['brainfuck', 'fish', 'j', 'hexagony'].includes(lang)
                    ? [] : extensions.bracketMatching,

                // These languages shouldn't wrap lines.
                ['assembly', 'fish', 'hexagony'].includes(lang)
                    ? [] : EditorView.lineWrapping,
            ],
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
    window.history.replaceState(null, '', '#' + select.value);
    switchLang();
}

switchLang();

// Run Code
const runCode = document.querySelector('#run a').onclick = async () => {
    document.querySelectorAll('canvas').forEach(e => e.remove());

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
    statusH2.innerText      = data.Pass ? 'Pass 😀'      : 'Fail ☹️';

    document.querySelector('#arguments').replaceChildren(...data.Argv.map(arg =>
        <span>{arg}</span>
    ));

    diffContent.innerHTML = diffHTML(hole, data.Exp, data.Out, data.Argv);

    document.querySelector('#errors').innerHTML   = data.Err;
    document.querySelector('#expected').innerText = data.Exp;
    document.querySelector('#output').innerText   = data.Out;

    // 3rd party integrations.
    thirdParty.replaceChildren( lang == 'hexagony' ? <a target="_blank" href={
        '//hexagony.net#lz' + LZString.compressToBase64(JSON.stringify({
            code, input: data.Argv.join('\0') + '\0', inputMode: 'raw' }))
    }>Run on Hexagony.net</a> : '' );

    if (hole == 'julia-set')
        document.querySelector('main').append(pbm(data.Exp), pbm(data.Out) ?? []);

    status.style.display = 'grid';

    update();
};

onkeydown = e => (e.ctrlKey || e.metaKey) && e.key == 'Enter' ? runCode() : undefined;

// Adapted from https://codegolf.stackexchange.com/a/119563
const ord = i => [, 'st', 'nd', 'rd'][i % 100 >> 3 ^ 1 && i % 10] || 'th';
