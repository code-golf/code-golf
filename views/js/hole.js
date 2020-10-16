const chars          = document.querySelector('#chars');
const details        = document.querySelector('#details');
const editor         = document.querySelector('#editor');
const hole           = decodeURI(location.pathname.slice(1));
const keymap         = JSON.parse(document.querySelector('#keymap').innerText);
const langs          = JSON.parse(document.querySelector('#langs').innerText);
const picker         = document.querySelector('#picker');
const scorings       = ['Bytes', 'Chars'];
const solutionPicker = document.querySelector('#solutionPicker');
const solutions      = JSON.parse(document.querySelector('#solutions').innerText);
const status         = document.querySelector('#status');
const table          = document.querySelector('#scores');

let lang;
let solution = Math.max(scorings.indexOf(localStorage.getItem('solution')), 0);
let scoring = Math.max(scorings.indexOf(localStorage.getItem('scoring')), 0);
let setCodeForLangAndSolution;
let latestSubmissionID = 0;

// Assume the user is logged in by default. At this point, it doesn't matter if the user is actually logged in,
// because this is only used to prevent auto-saving submitted solutions for logged in users to avoid excessive
// prompting to restore auto-saved solutions. The solutions dictionaries will initially be empty for users who are not
// logged in, so the loggedIn state will not be used. By the time they are non-empty, the loggedIn state will have been
// updated.
let loggedIn = true;

function getAutoSaveKey(lang, solution) {
    return `code_${hole}_${lang}_${solution}`;
}

function getSolutionCode(lang, solution) {
    return lang in solutions[solution] ? solutions[solution][lang] : '';
}

function getOtherScoring(value) {
    return 1 - value;
}

function setSolution(value) {
    localStorage.setItem('solution', scorings[solution = value]);
}

function setScoring(value) {
    localStorage.setItem('scoring', scorings[scoring = value]);
}

function formatScore(value) {
    return value ? value.toLocaleString('en') : '&nbsp';
}

onload = () => {
    // Lock the editor's height in so we scroll.
    editor.style.height = `${editor.offsetHeight}px`;

    const vimMode = keymap == 'vim';
    const cm      = new CodeMirror(editor, {
        autofocus:    true,
        indentUnit:   1,
        lineNumbers:  true,
        lineWrapping: true,
        smartIndent:  false,
        vimMode,
    });

    cm.on('change', () => {
        const code = cm.getValue();
        let infoText = '';
        for (let i = 0; i < scorings.length; i++) {
            if (i)
                infoText += ', ';
            infoText += `${getScoring(code, i).toLocaleString('en')} ${scorings[i].toLowerCase()}`;
        }
        chars.innerText = infoText;

        // Avoid future conflicts by only storing code locally that's different from the server's copy.
        const serverCode = getSolutionCode(lang, solution);

        const key = getAutoSaveKey(lang, solution);
        if (code && (code != serverCode || !loggedIn))
            localStorage.setItem(key, code);
        else
            localStorage.removeItem(key);
    });

    details.ontoggle = () =>
        document.cookie = 'hide-details=' + (details.open ? ';Max-Age=0' : '');

    setCodeForLangAndSolution = () => {
        const autoSaveCode = localStorage.getItem(getAutoSaveKey(lang, solution)) || '';
        const code = getSolutionCode(lang, solution) || autoSaveCode;

        cm.setOption('matchBrackets', lang != 'brainfuck' && lang != 'j');
        cm.setOption('mode', {name: 'text/x-' + lang, startOpen: true});
        cm.setValue(code);

        refreshScores();

        if (autoSaveCode && code != autoSaveCode &&
            confirm('Your local copy of the code is different than the remote one. Do you want to restore the local version?'))
            cm.setValue(autoSaveCode);

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

        if (res.status != 200) {
            alert('Error ' + res.status);
            return;
        }

        const data = await res.json();
        loggedIn = data.LoggedIn;

        if (submissionID != latestSubmissionID)
            return;

        if (data.Pass) {
            for (let i = 0; i < scorings.length; i++) {
                const solutionCode = getSolutionCode(codeLang, i);
                if (!solutionCode || getScoring(code, i) <= getScoring(solutionCode, i)) {
                    solutions[i][codeLang] = code;

                    // Don't need to keep solution in local storage because it's stored on the site.
                    // This prevents conflicts when the solution is improved on another browser.
                    if (loggedIn && localStorage.getItem(getAutoSaveKey(codeLang, i)) == code)
                        localStorage.removeItem(getAutoSaveKey(codeLang, i));
                }
            }
        }

        for (let i = 0; i < scorings.length; i++) {
            if (loggedIn) {
                // If the auto-saved code matches the other solution, remove it to avoid prompting the user to restore it.
                const autoSaveCode = localStorage.getItem(getAutoSaveKey(codeLang, i));
                for (let j = 0; j < scorings.length; j++) {
                    if (getSolutionCode(codeLang, j) == autoSaveCode)
                        localStorage.removeItem(getAutoSaveKey(codeLang, i));
                }
            }
            else {
                // Autosave the best solution for each scoring metric.
                localStorage.setItem(getAutoSaveKey(codeLang, i), getSolutionCode(codeLang, i));
            }
        }

        // Automatically switch to the solution whose code matches the current code after a new solution is submitted.
        // Don't change scoring. refreshScores will update the solution picker.
        if (data.Pass && getSolutionCode(codeLang, solution) != code && getSolutionCode(codeLang, getOtherScoring(solution)) == code)
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

    // Allow vim users to run code with :w or :write
    if (vimMode)
        CodeMirror.Vim.defineEx('write', 'w', submit);
};

async function refreshScores() {
    picker.innerHTML = '';
    picker.open = false;

    for (const l of langs) {
        let name = l.name;

        if (getSolutionCode(l.id, 0)) {
            name += ' <sup>';

            let lastValue = 0;
            for (let i = 0; i < scorings.length; i++) {
                const value = getScoring(getSolutionCode(l.id, i), i);
                if (lastValue != value) {
                    if (lastValue)
                        name += '/';
                    name += value.toLocaleString('en');
                }
                lastValue = value;
            }

            name += '</sup>';
        }

        picker.innerHTML += l.id == lang
            ? `<a>${name}</a>` : `<a href=#${l.id}>${name}</a>`;
    }

    while (solutionPicker.firstChild)
        solutionPicker.removeChild(solutionPicker.firstChild);

    const code0 = getSolutionCode(lang, 0);
    const code1 = getSolutionCode(lang, 1);
    const autoSave0 = localStorage.getItem(getAutoSaveKey(lang, 0));
    const autoSave1 = localStorage.getItem(getAutoSaveKey(lang, 1));

    // Only show the solution picker when both solutions are actually used.
    if (code0 && code1 && code0 != code1 || autoSave0 && autoSave1 && autoSave0 != autoSave1 ||
        (solution == 0 && code0 && autoSave1 && code0 != autoSave1) ||
        (solution == 1 && autoSave0 && code1 && autoSave0 != code1)) {
        for (let i = 0; i < scorings.length; i++) {
            let name = `Fewest ${scorings[i]}`;

            const solutionCode = getSolutionCode(lang, i);
            if (solutionCode) {
                name += ` <sup>${getScoring(solutionCode, i).toLocaleString('en')}</sup>`;
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

    const url = `/scores/${hole}/${lang}/${scorings[scoring].toLowerCase()}`;
    const res = await fetch(`${url}/mini`);

    const scores = res.ok ? await res.json() : [];
    let html     = `<thead><tr><th colspan=4>`;

    for (let i = 0; i < scorings.length; i++)
        html += `<a class="scoringPicker${scoring != i ? ' inactive' : ''}" id=${scorings[i]}>${scorings[i]}</a>`;

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
                <td class=right><span${scorings[scoring] != "Bytes" ? ' class=inactive' : ''}`;

            if (s.bytes)
                html += ` data-tooltip="Bytes solution is ${formatScore(s.bytes)} bytes, ${formatScore(s.bytes_chars)} chars."`;

            html += `>${formatScore(s.bytes)}</span>
                <td class=right><span${scorings[scoring] != "Chars" ? ' class=inactive' : ''}`;

            if (s.chars)
                html +=
                    ` data-tooltip="Chars solution is ${formatScore(s.chars_bytes)} bytes, ${formatScore(s.chars)} chars."`;

            html += `>${formatScore(s.chars)}</span>`;
        }
        else
            html += `<tr><td colspan=4>&nbsp`;
    }

    table.innerHTML = html;

    const otherScoring = getOtherScoring(scoring);
    const switchScoring = table.querySelector(`#${scorings[otherScoring]}`);
    switchScoring.href = 'javascript:void(0)';
    switchScoring.onclick = () => { setScoring(otherScoring); refreshScores(); };
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

function getScoring(str, index) {
    return scorings[index] == 'Bytes' ? new TextEncoder().encode(str).length : strlen(str);
}
