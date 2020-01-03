/* include vendor/codemirror.js            */
/* include vendor/codemirror-bash.js       */
/* include vendor/codemirror-clike.js      */
/* include vendor/codemirror-haskell.js    */
/* include vendor/codemirror-htmlmixed.js  */
/* include vendor/codemirror-javascript.js */
/* include vendor/codemirror-julia.js      */
/* include vendor/codemirror-lisp.js       */
/* include vendor/codemirror-lua.js        */
/* include vendor/codemirror-nim.js        */
/* include vendor/codemirror-perl.js       */
/* include vendor/codemirror-php.js        */
/* include vendor/codemirror-python.js     */
/* include vendor/codemirror-raku.js       */
/* include vendor/codemirror-ruby.js       */
/* include vendor/codemirror-xml.js        */

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

onload = () => {
    // Lock the editor's height in so we scroll.
    editor.style.height = `${editor.offsetHeight}px`;

    const cm = new CodeMirror(editor, {autofocus: true, lineNumbers: true, lineWrapping: true});

    cm.on('change', () => {
        const val = cm.getValue();
        const len = strlen(val);

        chars.innerText = `${len.toLocaleString('en')} character${len - 1 ? 's' : ''}`;
    });

    details.ontoggle = () =>
        document.cookie = 'hide-details=' + (details.open ? ';Max-Age=0' : '');

    (onhashchange = () => {
        lang = location.hash.slice(1);

        // Kick 'em to Python if we don't know the chosen language.
        if (!langs.find(l => l.id == lang))
            location.hash = lang = 'python';

        cm.setOption('mode', {name: lang == 'c' ? 'clike' : lang, startOpen: true});
        cm.setValue(lang in solutions ? solutions[lang] : '');

        picker.innerHTML = '';
        picker.open = false;

        for (const l of langs) {
            let name = l.name;

            if (l.id in solutions)
                name += ` <sup>${strlen(solutions[l.id]).toLocaleString('en')}</sup>`;

            picker.innerHTML += l.id == lang
                ? `<a>${name}</a>` : `<a href=#${l.id}>${name}</a>`;
        }

        refreshScores();
    })();

    const submit = document.querySelector('#run a').onclick = async () => {
        const res  = await fetch('/solution', {
            method: 'POST',
            body: JSON.stringify({
                Code: cm.getValue(),
                Hole: hole,
                Lang: lang,
            }),
        });

        const data = await res.json();
        const pass = data.Exp === data.Out && data.Out !== '';

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
