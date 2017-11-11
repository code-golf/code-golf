/* include codemirror.js            */
/* include codemirror-bash.js       */
/* include codemirror-clike.js      */
/* include codemirror-javascript.js */
/* include codemirror-perl6.js      */
/* include codemirror-perl.js       */
/* include codemirror-php.js        */
/* include codemirror-python.js     */
/* include codemirror-ruby.js       */

onload = function() {
    let activeEditor;
    let article = document.getElementsByTagName('article')[0];
    let editors = [];

    for (let lang of ['bash', 'javascript', 'perl', 'perl6', 'php', 'python', 'ruby']) {
        let code = article.dataset.hasOwnProperty(lang) ? article.dataset[lang] : '';

        let editor = CodeMirror(article, {
            lineNumbers: true,
            lineWrapping: true,
            mode: { name: lang, startOpen: true },
            value: code,
        });

        let span = document.querySelector('[href="#' + lang + '"] span');

        let callback = function(editor) {
            let len = [...editor.getValue()].length;

            span.innerText = len === 1 ? '1 char'
                           : len       ? len + ' chars'
                           :             'not tried';
        };

        callback(editor);

        editor.on('change', callback);

        editors.push(editor);
    }

    ( onhashchange = function() {
        // Kick 'em to Perl 6 if we don't know the chosen language.
        if (!/^#(?:bash|javascript|perl6?|php|python|ruby)$/.exec(location.hash))
            location.hash = 'perl6';

        let lang = location.hash.slice(1);

        for (let editor of editors)
            if (editor.options.mode.name === lang)
                ( activeEditor = editor ).display.wrapper.style.display = '';
            else
                editor.display.wrapper.style.display = 'none';

        for (let tab of document.getElementsByClassName('tab'))
            tab.classList.toggle('on', tab.href === location.href);
    } )();

    document.getElementsByTagName('input')[0].onclick = function() {
        fetch('/solution', {
            credentials: 'include',
            method: 'POST',
            body: JSON.stringify({
                Code: activeEditor.getValue(),
                Hole: decodeURI(location.pathname.slice(1)),
                Lang: activeEditor.options.mode.name,
            }),
        }).then( res => res.json() ).then( data => {
            document.getElementById('status').classList
                .toggle('pass', data.Exp === data.Out && data.Out !== '');

            console.log(data.Diff);

            for (let prop in data) {
                if (prop === 'Diff')
                    continue;

                let pre = document.getElementById(prop);

                // Err can be ANSI coloured via HTML.
                if (prop == Err)
                    pre.innerHTML = data[prop];
                else
                    pre.innerText = data[prop];

                // Only show Arg & Err if they contain something.
                if (prop === 'Arg' || prop === 'Err')
                    pre.style.display = pre.previousSibling.style.display
                        = data[prop] ? '' : 'none';
            }

            document.getElementById('status').style.display = 'block';
        });
    };
};
