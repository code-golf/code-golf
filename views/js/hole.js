const chars          = document.querySelector('#chars');
const details        = document.querySelector('#details');
const editor         = document.querySelector('#editor');
const hole           = decodeURI(location.pathname.slice(1));
const langs          = JSON.parse(document.querySelector('#langs').innerText);
const picker         = document.querySelector('#picker');
const solutionPicker = document.querySelector('#solutionPicker');
const scorings       = JSON.parse(document.querySelector('#scorings').innerText);
const solutions      = JSON.parse(document.querySelector('#solutions').innerText);
const status         = document.querySelector('#status');
const table          = document.querySelector('#scores');
const beta           = scorings.length > 1;

let lang;
let solution = beta ? Math.max(scorings.indexOf(localStorage.getItem('solution')), 0) : 0;
let scoring = beta ? Math.max(scorings.indexOf(localStorage.getItem('scoring')), 0) : 0;
let setCodeForLangAndSolution;
let latestSubmissionID = 0;

function getOtherScoring(value) {
    return beta ? 1 - value : value;
}

function setSolution(value) {
    localStorage.setItem('solution', scorings[solution = value]);
}

function setScoring(value) {
    localStorage.setItem('scoring', scorings[scoring = value]);
}

onload = () => {
    // Lock the editor's height in so we scroll.
    editor.style.height = `${editor.offsetHeight}px`;

    const cm = new CodeMirror(editor, {
        autofocus:    true,
        indentUnit:   1,
        lineNumbers:  true,
        lineWrapping: true,
        smartIndent:  false,
    });

    cm.on('change', () => {
        const code = cm.getValue();
        const len = strlen(code);
        chars.innerText = `${len.toLocaleString('en')} character${len - 1 ? 's' : ''}`;
        if (beta) {
            const bytes = utf8ByteCount(code);
            if (bytes != len)
                chars.innerText += `, ${bytes.toLocaleString('en')} byte${bytes - 1 ? 's' : ''}`;
        }

        // Avoid future conflicts by only storing code locally that's different from the server's copy.
        const serverCode = lang in solutions[solution] ? solutions[solution][lang] : '';

        const key = `code_${hole}_${lang}_${solution}`;
        if (code && code != serverCode)
            localStorage.setItem(key, code);
        else
            localStorage.removeItem(key);
    });

    details.ontoggle = () =>
        document.cookie = 'hide-details=' + (details.open ? ';Max-Age=0' : '');

    setCodeForLangAndSolution = () => {
        const code = lang in solutions[solution] ? solutions[solution][lang] : '';
        const previousCode = localStorage.getItem(`code_${hole}_${lang}_${solution}`);

        cm.setOption('matchBrackets', lang != 'brainfuck' && lang != 'j');
        cm.setOption('mode', {name: 'text/x-' + lang, startOpen: true});
        cm.setValue(code);

        refreshScores();

        if (previousCode && code != previousCode && (!code ||
            confirm('Your local copy of the code is different than the remote one. Do you want to restore the local version?')))
            cm.setValue(previousCode);

        for (let info of document.querySelectorAll('main .info'))
            info.style.display = info.classList.contains(lang) ? 'block' : '';
    };

    (onhashchange = () => {
        lang = location.hash.slice(1) || localStorage.getItem('lang');

        // Kick 'em to Python if we don't know the chosen language.
        if (!langs.find(l => l.id == lang))
            lang = 'python';

        localStorage.setItem('lang', lang);

        history.replaceState(null, '', '#' + lang);

        setCodeForLangAndSolution();
    })();

    const submit = document.querySelector('#run a').onclick = async () => {
        document.querySelector('h2').innerText = '…';
        status.className = 'grey';

        const code = cm.getValue();
        const codeLang = lang;
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

        for (let i = 0; i < scorings.length; i++) {
            const lengthFunc = i ? utf8ByteCount : strlen;
            if (data.Pass && (!(codeLang in solutions[i]) || lengthFunc(code) <= lengthFunc(solutions[i][codeLang]))) {
                solutions[i][codeLang] = code;

                // Don't need to keep solution in local storage because it's stored on the site.
                // This prevents conflicts when the solution is improved on another browser.
                if (data.LoggedIn)
                    localStorage.removeItem(`code_${hole}_${codeLang}_${i}`);
            }
        }

        // Automatically switch to the solution whose code matches the current code after a new solution is submitted.
        // Don't change scoring. refreshScores will update the solution picker.
        if (beta && data.Pass && solutions[solution][lang] != code && solutions[getOtherScoring(solution)][lang] == code)
            setSolution(getOtherScoring(solution));

        document.querySelector('h2').innerText
            = data.Pass ? 'Pass 😀' : 'Fail ☹️';

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
        if (data.Err && !data.Pass) {
            document.querySelector('#err').style.display = 'block';
            document.querySelector('#err div').innerHTML = data.Err.replace(/\n/g, '<br>');
        }
        else
            document.querySelector('#err').style.display = '';

        // Always show exp & out.
        document.querySelector('#exp div').innerText = data.Exp;
        document.querySelector('#out div').innerText = data.Out;

        status.className = data.Pass ? 'green' : 'red';
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

        if (l.id in solutions[0]) {
            const chars = strlen(solutions[0][l.id])
            let sameCount = true;

            if (beta) {
                const bytes = utf8ByteCount(solutions[1][l.id])
                if (chars != bytes) {
                    name += ` <sup>${chars.toLocaleString('en')}/${bytes.toLocaleString('en')}</sup>`;
                    sameCount = false;
                }
            }

            if (sameCount)
                name += ` <sup>${chars.toLocaleString('en')}</sup>`;
        }

        picker.innerHTML += l.id == lang
            ? `<a>${name}</a>` : `<a href=#${l.id}>${name}</a>`;
    }

    if (beta) {
        while (solutionPicker.firstChild)
            solutionPicker.removeChild(solutionPicker.firstChild);

        // Only show the solution picker when both solutions are actually used.
        if (lang in solutions[0] && lang in solutions[1] && solutions[0][lang] != solutions[1][lang]) {
            for (let i = 0; i < scorings.length; i++) {
                let name = `Fewest ${scorings[i]}`;

                if (lang in solutions[i]) {
                    name += ` <sup>${[strlen,utf8ByteCount][i](solutions[i][lang]).toLocaleString('en')}</sup>`;
                }

                const child = document.createElement('a');
                child.innerHTML = name;
                if (i != solution) {
                    child.href = 'javascript:void(0)';
                    child.onclick = () => {
                        setSolution(i);
                        setCodeForLangAndSolution();
                    };
                }
                solutionPicker.appendChild(child);
            }
        }
    }

    const url    = `/scores/${hole}/${lang}`;
    const scores = await (await fetch(`${url}/${scoring ? 'minibytes' : 'mini'}`)).json();
    let html     = `<thead><tr>`;

    if (beta) {
        html += `<thead><tr><th colspan=4>`;
        for (let i = 0; i < scorings.length; i++)
            html += `<a class="scoringPicker${scoring != i ? ' inactive' : ''}" id=${scorings[i]}>${scorings[i]}</a>`;
    }
    else
        html += `<th colspan=3>Scores`;

    html += `<a href=${url} id=all>all</a><tbody>`;

    // Ordinal from https://codegolf.stackexchange.com/a/119563
    for (let i = 0; i < 7; i++) {
        const s = scores[i];

        if (s) {
            html += `<tr ${s.me ? 'class=me' : ''}>
                <td>${s.rank}<sup>${[, 'st', 'nd', 'rd'][s.rank % 100 >> 3 ^ 1 && s.rank % 10] || 'th'}</sup>
                <td><a href=/golfers/${s.login}>
                    <img src="//avatars.githubusercontent.com/${s.login}?s=24">
                    <span>${s.login}</span>
                </a>
                <td class=right><span${scoring != 0 ? ' class=inactive' : ''}>${s.chars ? s.chars.toLocaleString('en') : '&nbsp'}</span>`;

            if (beta)
                html += `<td class=right><span${scoring != 1 ? ' class=inactive' : ''}>${s.bytes ? s.bytes.toLocaleString('en') : '&nbsp'}</span>`;
        }
        else
            html += `<tr><td colspan=${beta ? 4 : 3}>&nbsp`;
    }

    table.innerHTML = html;

    if (beta) {
        const otherScoring = getOtherScoring(scoring);
        const switchScoring = table.querySelector(`#${scorings[otherScoring]}`);
        switchScoring.href = 'javascript:void(0)';
        switchScoring.onclick = () => { setScoring(otherScoring); refreshScores(); };
    }
}

function bytesForCodePoint(codePoint) {
    if ((codePoint & 0xFFFFFF80) == 0) {
        return 1;
    }
    if ((codePoint & 0xFFFFF800) == 0) {
        return 2;
    }
    else if ((codePoint & 0xFFFF0000) == 0) {
        return 3;
    }
    return 4;
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

// Adapted from https://mths.be/punycode
function utf8ByteCount(str) {
    let i = 0, byteCount = 0;

    while (i < str.length) {
        const value = str.charCodeAt(i++);

        if (value >= 0xD800 && value <= 0xDBFF && i < str.length) {
            // It's a high surrogate, and there is a next character.
            const extra = str.charCodeAt(i++);

            // Low surrogate.
            if ((extra & 0xFC00) == 0xDC00) {
                byteCount += bytesForCodePoint(((value & 0x3FF) << 10) + (extra & 0x3FF) + 0x10000);
            } else {
                // It's an unmatched surrogate; only append this code unit, in case the
                // next code unit is the high surrogate of a surrogate pair.
                byteCount += bytesForCodePoint(value);
                i--;
            }
        } else {
            byteCount += bytesForCodePoint(value);
        }
    }

    return byteCount;
}
