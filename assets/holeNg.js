/* include prism.js */

const hole = decodeURI(location.pathname.slice(4));

let activeTextarea, activeStrokes;

function calcStrokes(textarea) {
    const len = [...textarea.value].length;

    return len === 1 ? '1 char'
         : len       ? len + ' chars'
         :             'not tried';
}

function render() {
    let code = activeTextarea.nextSibling.children[0];

    code.textContent = activeTextarea.value + "\n";

    activeStrokes.textContent = calcStrokes(activeTextarea);

    Prism.highlightElement(code);
}

onload = function() {
    const status = document.querySelector('#status');

    let divs = document.querySelectorAll('.code');

    for (let div of divs) {
        let code     = document.createElement('code');
        let pre      = document.createElement('pre');
        let textarea = document.createElement('textarea');

        code.classList.add('lang-' + div.dataset.lang);
        pre.classList.add('line-numbers');
        textarea.setAttribute('autocapitalize', 'off');
        textarea.setAttribute('autocomplete', 'off');
        textarea.setAttribute('autocorrect', 'off');
        textarea.setAttribute('spellcheck', 'false');
        textarea.value = div.textContent;

        document.querySelector('[href="#' + div.dataset.lang + '"] div:nth-child(2)')
                .textContent = calcStrokes(textarea);

        // Fixing iOS "drunk-text" issue
        if(/iPad|iPhone|iPod/.test(navigator.platform))
            code.style.paddingLeft = '3px';

        div.innerHTML = '';
        div.appendChild(textarea);
        div.appendChild(pre);
        pre.appendChild(code);

        textarea.oninput   = render;
        textarea.onkeydown = function(e) {
            // If tab is pressed, insert four spaces.
            if(e.keyCode === 9){
                e.preventDefault();

                let start = this.selectionStart;

                this.value = this.value.substring(0, start)
                           + '    '
                           + this.value.substring(start, this.value.length);

                this.selectionStart = start + 4;
                this.selectionEnd   = start + 4;

                render();
            }
        };

        textarea.onscroll = function(){
            pre.style.top = '-' + Math.floor(this.scrollTop) + 'px';
        };
    }

    // Get the tab from either the URL, local storage, or the latest code.
    location.hash = location.hash.slice(1) || localStorage.getItem(hole);

    ( onhashchange = function() {
        // Kick 'em to Perl 6 if we don't know the chosen language.
        if (!/^#(?:bash|haskell|javascript|lisp|lua|perl6?|php|python|ruby)$/.exec(location.hash))
            location.hash = 'perl6';

        let lang = location.hash.slice(1);

        localStorage.setItem(hole, lang);

        for (let div of divs)
            if (div.dataset.lang === lang) {
                div.style.display = 'block';
                activeTextarea = div.children[0];
            }
            else
                div.style.display = '';

        for (let tab of document.querySelectorAll('#tabs a'))
            if (tab.href === location.href) {
                tab.classList.add('on');
                activeStrokes = tab.children[1];
            }
            else
                tab.classList.remove('on');

        render();
    } )();

    onresize = render;

    document.querySelector('button').onclick = function() {
        status.style.display = '';
        this.classList.add('on');

        fetch('/solution', {
            credentials: 'include',
            method: 'POST',
            body: JSON.stringify({
                Code: activeTextarea.value,
                Hole: hole,
                Lang: activeTextarea.parentNode.dataset.lang,
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
