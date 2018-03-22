/* include codemirror.js            */
/* include codemirror-bash.js       */
/* include codemirror-clike.js      */
/* include codemirror-javascript.js */
/* include codemirror-perl6.js      */
/* include codemirror-perl.js       */
/* include codemirror-php.js        */
/* include codemirror-python.js     */
/* include codemirror-ruby.js       */

const hole = decodeURI(location.pathname.slice(1));

onload = function() {
    const main   = document.querySelector('main');
    const data   = main.dataset;
    const status = document.querySelector('#status');

    let activeEditor;
    let editors = [];

    for (let lang of ['bash', 'javascript', 'lisp', 'lua', 'perl', 'perl6', 'php', 'python', 'ruby']) {
        let editor = CodeMirror(main, {
            lineNumbers: true,
            lineWrapping: true,
            mode: { name: lang, startOpen: true },
            value: data.hasOwnProperty(lang) ? data[lang] : '',
        });

        let div = document.querySelector('[href="#' + lang + '"] div:nth-child(2)');

        let callback = function(editor) {
            let len = [...editor.getValue()].length;

            div.innerText = len === 1 ? '1 char'
                          : len       ? len + ' chars'
                          :             'not tried';
        };

        callback(editor);

        editor.on('change', callback);

        editors.push(editor);
    }

    // Get the tab from either the URL, local storage, or the latest code.
    location.hash = location.hash.slice(1) || localStorage.getItem(hole) || data.lang;

    ( onhashchange = function() {
        // Kick 'em to Perl 6 if we don't know the chosen language.
        if (!/^#(?:bash|javascript|lisp|lua|perl6?|php|python|ruby)$/.exec(location.hash))
            location.hash = 'perl6';

        let lang = location.hash.slice(1);

        localStorage.setItem(hole, lang);

        for (let editor of editors)
            if (editor.options.mode.name === lang)
                ( activeEditor = editor ).display.wrapper.style.display = '';
            else
                editor.display.wrapper.style.display = 'none';

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
            status.classList.toggle('pass', data.Exp === data.Out && data.Out !== '');

            console.log(data.Diff);

            for (let prop in data) {
                if (prop === 'Diff')
                    continue;

                let pre = document.getElementById(prop);

                // Err can be ANSI coloured via HTML.
                if (prop === 'Err')
                    pre.innerHTML = data[prop];
                else
                    pre.innerText = data[prop];

                // Only show Arg & Err if they contain something.
                if (prop === 'Arg' || prop === 'Err')
                    pre.style.display = pre.previousSibling.style.display
                        = data[prop] ? '' : 'none';
            }

            status.style.display = 'block';
            this.classList.remove('on');
        });
    };
};
