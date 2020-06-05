/* include vendor/codemirror.js            */
/* include vendor/codemirror-clike.js      */
/* include vendor/codemirror-htmlmixed.js  */
/* include vendor/codemirror-simple.js     */
/* include vendor/codemirror-xml.js        */

/* include vendor/codemirror-bash.js       */
/* include vendor/codemirror-brainfuck.js  */
/* include vendor/codemirror-fortran.js    */
/* include vendor/codemirror-go.js         */
/* include vendor/codemirror-haskell.js    */
/* include vendor/codemirror-javascript.js */
/* include vendor/codemirror-julia.js      */
/* include vendor/codemirror-lisp.js       */
/* include vendor/codemirror-lua.js        */
/* include vendor/codemirror-mllike.js     */
/* include vendor/codemirror-nim.js        */
/* include vendor/codemirror-perl.js       */
/* include vendor/codemirror-php.js        */
/* include vendor/codemirror-powershell.js */
/* include vendor/codemirror-python.js     */
/* include vendor/codemirror-raku.js       */
/* include vendor/codemirror-ruby.js       */
/* include vendor/codemirror-rust.js       */
/* include vendor/codemirror-swift.js      */

const chars     = document.querySelector('#chars');
const details   = document.querySelector('#details');
const editor    = document.querySelector('#editor');
const hole      = decodeURI(location.pathname.slice(1));
const langs     = JSON.parse(document.querySelector('#langs').innerText);
const picker    = document.querySelector('#picker');
const solutions = JSON.parse(document.querySelector('#solutions').innerText);
const status    = document.querySelector('#status');
const table     = document.querySelector('.scores');

let lang;
let latestSubmissionID = 0;

onload = () => {
    if (hole != 'quine')
        document.querySelectorAll('.quine').forEach(e => e.remove());

    // Lock the editor's height in so we scroll.
    editor.style.height = `${editor.offsetHeight}px`;

    const cm = new CodeMirror(editor, {autofocus: true, lineNumbers: true, lineWrapping: true, smartIndent: false});

    cm.on('change', () => {
        const code = cm.getValue();
        const len = strlen(code);

        chars.innerText = `${len.toLocaleString('en')} character${len - 1 ? 's' : ''}`;
        localStorage.setItem('code_' + hole + '_' + lang, code);
    });

    details.ontoggle = () =>
        document.cookie = 'hide-details=' + (details.open ? ';Max-Age=0' : '');

    (onhashchange = () => {
        lang = location.hash.slice(1) || localStorage.getItem('lang');

        // Kick 'em to Python if we don't know the chosen language.
        if (!langs.find(l => l.id == lang))
            location.hash = lang = 'python';

        const code = lang in solutions ? solutions[lang] : '';
        const previousCode = localStorage.getItem('code_' + hole + '_' + lang);

        cm.setOption('mode', {name: 'text/x-' + lang, startOpen: true});
        cm.setValue(code);

        localStorage.setItem('lang', location.hash = lang);

        refreshScores();

        if (previousCode && code != previousCode && (!code ||
            confirm('Your local copy of the code is different than the remote one. Do you want to restore the local version?')))
            cm.setValue(previousCode);

        for (let info of document.querySelectorAll('.info'))
            info.style.display = info.classList.contains(lang) ? 'block' : '';
    })();

    const submit = document.querySelector('#run a').onclick = async () => {
        document.querySelector('h2').innerText = '…';
        status.className = 'grey';

        const code = cm.getValue();
        const submissionID = ++latestSubmissionID;

        const res  = await fetch('/solution', {
            method: 'POST',
            body: JSON.stringify({
                Code: code,
                Hole: hole,
                Lang: lang,
            }),
        });

        const data = await res.json();
        if (submissionID != latestSubmissionID)
            return;

        const pass = data.Exp === data.Out && data.Out !== '';

        if (pass && (!(lang in solutions) || strlen(code) <= strlen(solutions[lang])))
            solutions[lang] = code;

        document.querySelector('h2').innerText
            = pass ? 'Pass 😊️' : 'Fail ☹️';

        // Show args if we have 'em.
        if (data.Argv) {
            document.querySelector('#arg').style.display = 'block';
            const argDiv = document.querySelector('#arg div');
            // Remove all arg spans
            while (argDiv.firstChild) {
                argDiv.removeChild(argDiv.firstChild);
            }
            // Add a span for each arg
            for (const arg of data.Argv) {
                argDiv.appendChild(document.createElement('span'));
                argDiv.lastChild.innerText = arg;
                argDiv.appendChild(document.createTextNode(' '));
            }
        }
        else
            document.querySelector('#arg').style.display = '';

        // Show err if we have some and we're not passing.
        if (data.Err && !pass) {
            document.querySelector('#err').style.display = 'block';
            document.querySelector('#err div').innerHTML = data.Err.replace(/\n/g, '<br>');
        }
        else
            document.querySelector('#err').style.display = '';

        // Always show exp & out.
        document.querySelector('#exp div').innerText = data.Exp;
        document.querySelector('#out div').innerText = data.Out;

        status.className = pass ? 'green' : 'red';
        status.style.display = 'block';

        refreshScores();
    };

    onkeydown = e => e.ctrlKey && e.key == 'Enter' ? submit() : undefined;
};

async function refreshScores() {
    picker.innerHTML = '';
    picker.open = false;

    for (const l of langs) {
        let name = l.name;

        if (l.id in solutions)
            name += ` <sup>${strlen(solutions[l.id]).toLocaleString('en')}</sup>`;

        picker.innerHTML += l.id == lang
            ? `<a>${name}</a>` : `<a href=#${l.id}>${name}</a>`;
    }

    const url    = `/scores/${hole}/${lang}`;
    const scores = await (await fetch(`${url}/mini`)).json();

    let html = `<thead><tr><th colspan=3>Scores<a href=${url}>all</a><tbody>`;

    // Ordinal from https://codegolf.stackexchange.com/a/119563
    for (let i = 0; i < 7; i++) {
        const s = scores[i];

        html += s ? `<tr ${s.me ? 'class=me' : ''}>
            <td>${s.rank}<sup>${[, 'st', 'nd', 'rd'][s.rank % 100 >> 3 ^ 1 && s.rank % 10] || 'th'}</sup>
            <td><a href=/users/${s.login}>
                <img src="//avatars.githubusercontent.com/${s.login}?s=24">${s.login}
            </a>
            <td class=right>${s.strokes.toLocaleString('en')}` : '<tr><td colspan=3>&nbsp;';

        // If the user's solution is from before the powershell change (b79e88ca)
        if (lang == 'powershell' && hole != 'quine' && s.me && s.submitted < '2020-06-04T19:56:49') {
            const info = document.querySelector('.info.special');
            info.innerHTML = '<b>Write-Host</b> is no longer required for output.';
            info.style.display = 'block';
        }
    }

    table.innerHTML = html;
}

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
