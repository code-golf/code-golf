onload = function() {
    let activeEditor;
    let article = document.getElementsByTagName('article')[0];
    let editors = [];

    for (let lang of ['javascript', 'perl', 'perl6', 'php', 'python', 'ruby'])
        editors.push(CodeMirror(article, {
            lineNumbers: true,
            lineWrapping: true,
            mode: { name: lang, startOpen: true },
            value: article.dataset.hasOwnProperty(lang) ? article.dataset[lang] : '',
        }));

    ( onhashchange = function() {
        // Kick 'em to Perl 6 if we don't know the chosen language.
        if (!/^#(?:javascript|perl6?|php|python|ruby)$/.exec(location.hash))
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
                .toggle('pass', data.Exp === data.Out);

            for (let prop in data)
                document.getElementById(prop).value = data[prop];

            document.getElementById('status').style.display = 'block';
        });
    };
};
