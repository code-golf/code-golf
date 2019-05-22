/* include codemirror.js            */
/* include codemirror-bash.js       */
/* include codemirror-clike.js      */
/* include codemirror-haskell.js    */
/* include codemirror-htmlmixed.js  */
/* include codemirror-javascript.js */
/* include codemirror-julia.js      */
/* include codemirror-lisp.js       */
/* include codemirror-lua.js        */
/* include codemirror-nim.js        */
/* include codemirror-perl.js       */
/* include codemirror-perl6.js      */
/* include codemirror-php.js        */
/* include codemirror-python.js     */
/* include codemirror-ruby.js       */
/* include codemirror-xml.js        */

const hole = decodeURI(location.pathname.slice(1));

// Adapted from https://mths.be/punycode
function ucs2decode(string) {
    let chars = 0;
    let counter = 0;
    const length = string.length;
    while (counter < length) {
        const value = string.charCodeAt(counter++);
        if (value >= 0xD800 && value <= 0xDBFF && counter < length) {
            // It's a high surrogate, and there is a next character.
            const extra = string.charCodeAt(counter++);
            if ((extra & 0xFC00) == 0xDC00) { // Low surrogate.
                chars++;
            } else {
                // It's an unmatched surrogate; only append this code unit, in case the
                // next code unit is the high surrogate of a surrogate pair.
                chars++;
                counter--;
            }
        } else {
            chars++;
        }
    }
    return chars;
}

onload = function() {
    const data    = document.querySelector('main').dataset;
    const status  = document.querySelector('#status');
    const wrapper = document.querySelector('#wrapper');

    let activeEditor;
    let editors = [];

    for (let langName of [
        'Bash', 'Brainfuck', 'C', 'Haskell', 'J', 'JavaScript', 'Julia',
        'Lisp', 'Lua', 'Nim', 'Perl', 'Perl 6', 'PHP', 'Python', 'Ruby',
    ]) {
        let lang = langName.replace(/ /, '').toLowerCase();

        let editor = CodeMirror(wrapper, {
            lineNumbers: true,
            lineWrapping: true,
            mode: { name: lang === 'c' ? 'clike' : lang, startOpen: true },
            value: data.hasOwnProperty(lang) ? data[lang] : '',
        });

        let key = hole + '-' + lang;
        let tab = document.querySelector('[href="#' + lang + '"]');

        let callback = function(editor) {
            const val = editor.getValue();
            const len = ucs2decode(val);

            tab.innerText = len || '';

            if (len)
                localStorage.setItem(key, val);
            else
                localStorage.removeItem(key);
        };

        let localCode = localStorage.getItem(key);
        if (localCode != null && localCode != editor.getValue())
            if (confirm('Restore your local ' + langName + ' solution?'))
                editor.setValue(localCode);

        callback(editor);

        editor.on('change', callback);

        editors.push(editor);
    }

    // Get the tab from either the URL, local storage, or the latest code.
    location.hash = location.hash.slice(1) || localStorage.getItem(hole) || data.lang;

    ( onhashchange = function() {
        // Kick 'em to Perl 6 if we don't know the chosen language.
        if (!/^#(?:bash|brainfuck|c|haskell|j|javascript|julia|lisp|lua|nim|perl6?|php|python|ruby)$/.exec(location.hash))
            location.hash = 'perl6';

        let lang = location.hash.slice(1);

        localStorage.setItem(hole, lang);

        for (let editor of editors) {
            let mode = editor.options.mode.name;

            if (mode === 'clike')
                mode = 'c';

            if (mode === lang)
                ( activeEditor = editor ).display.wrapper.style.display = '';
            else
                editor.display.wrapper.style.display = 'none';
        }

        for (let info of document.querySelectorAll('.info'))
            info.style.display = info.classList.contains(lang) ? 'block' : '';

        for (let tab of document.querySelectorAll('#tabs a'))
            tab.classList.toggle('on', tab.href === location.href);
    } )();

    document.querySelector('button').onclick = async function() {
        status.style.display = 'none';
        this.classList.add('on');

        const mode = activeEditor.options.mode.name;
        const res  = await fetch('/solution', {
            credentials: 'include',
            method: 'POST',
            body: JSON.stringify({
                Code: activeEditor.getValue(),
                Hole: hole,
                Lang: mode === 'clike' ? 'c' : mode,
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

        status.classList.toggle('pass', pass);
        status.style.display = 'block';
        this.classList.remove('on');
    };

    const h1 = document.querySelector('h1');

    if (localStorage.getItem(hole+'-collapsed'))
        h1.className = 'collapsed';

    h1.onclick = function() {
        if (h1.classList.toggle('collapsed'))
            localStorage.setItem(hole+'-collapsed', 1);
        else
            localStorage.removeItem(hole+'-collapsed');
    };
};
