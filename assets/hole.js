/* include codemirror.js            */
/* include codemirror-bash.js       */
/* include codemirror-clike.js      */
/* include codemirror-haskell.js    */
/* include codemirror-htmlmixed.js  */
/* include codemirror-javascript.js */
/* include codemirror-lisp.js       */
/* include codemirror-lua.js        */
/* include codemirror-perl.js       */
/* include codemirror-perl6.js      */
/* include codemirror-php.js        */
/* include codemirror-python.js     */
/* include codemirror-ruby.js       */
/* include codemirror-xml.js        */

const hole = decodeURI(location.pathname.slice(1));

onload = function() {
    const data    = document.querySelector('main').dataset;
    const status  = document.querySelector('#status');
    const wrapper = document.querySelector('#wrapper');

    let activeEditor;
    let editors = [];

    for (let lang of ['bash', 'haskell', 'javascript', 'lisp', 'lua', 'perl', 'perl6', 'php', 'python', 'ruby']) {
        let editor = CodeMirror(wrapper, {
            lineNumbers: true,
            lineWrapping: true,
            mode: { name: lang, startOpen: true },
            value: data.hasOwnProperty(lang) ? data[lang] : '',
        });

        let tab = document.querySelector('[href="#' + lang + '"]');

        let callback = function(editor) {
            let len = [...editor.getValue()].length;

            tab.innerText = len === 1 ? '1 char'
                          : len       ? len + ' chars'
                          :             'not tried';

            localStorage.setItem("code_"+hole+"_"+lang, editor.getValue())
        };

        let previousCode = localStorage.getItem("code_"+hole+"_"+lang);
        if (previousCode && previousCode != editor.getValue())
            if (confirm("Your local copy of the code is different than the remote one. Do you want to restore the local version?"))
                editor.setValue(previousCode);

        callback(editor);

        editor.on('change', callback);

        editors.push(editor);
    }

    // Get the tab from either the URL, local storage, or the latest code.
    location.hash = location.hash.slice(1) || localStorage.getItem(hole) || data.lang;

    ( onhashchange = function() {
        // Kick 'em to Perl 6 if we don't know the chosen language.
        if (!/^#(?:bash|haskell|javascript|lisp|lua|perl6?|php|python|ruby)$/.exec(location.hash))
            location.hash = 'perl6';

        let lang = location.hash.slice(1);

        localStorage.setItem(hole, lang);

        for (let editor of editors)
            if (editor.options.mode.name === lang)
                ( activeEditor = editor ).display.wrapper.style.display = '';
            else
                editor.display.wrapper.style.display = 'none';

        for (let info of document.querySelectorAll('.info'))
            info.style.display = info.classList.contains(lang) ? 'block' : '';

        for (let tab of document.querySelectorAll('#tabs a'))
            tab.classList.toggle('on', tab.href === location.href);
    } )();

    document.querySelector('button').onclick = function() {
        status.style.display = 'none';
        this.classList.add('on');

        fetch('/solution', {
            credentials: 'include',
            method: 'POST',
            body: JSON.stringify({
                Code: activeEditor.getValue(),
                Hole: hole,
                Lang: activeEditor.options.mode.name,
            }),
        }).then( res => res.json() ).then( data => {
            const pass = data.Exp === data.Out && data.Out !== '';

            document.querySelector('h2').innerText
                = pass ? 'Pass 😊️' : 'Fail ☹️';

            // Show args if we have 'em.
            if (data.Argv) {
                document.querySelector('#arg').style.display = 'block';
                document.querySelector('#arg div').innerHTML
                    = '<span>' + data.Argv.join('</span> <span>') + '</span>';
            }
            else
                document.querySelector('#arg').style.display = '';

            // Show err if we have some and we're not passing.
            if (data.Err && !pass) {
                document.querySelector('#err').style.display = 'block';
                document.querySelector('#err div').innerHTML = data.Err;
            }
            else
                document.querySelector('#err').style.display = '';

            // Always show exp & out.
            document.querySelector('#exp div').innerText = data.Exp;
            document.querySelector('#out div').innerText = data.Out;

            status.classList.toggle('pass', pass);
            status.style.display = 'block';
            this.classList.remove('on');
        });
    };
};
