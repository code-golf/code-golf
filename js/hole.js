import CodeMirror                       from './_codemirror-legacy';
import                                       './_copy-as-json';
import { attachDiff }                   from './_diff';
import { byteLen, charLen, comma, ord } from './_util';

const chars              = document.querySelector('#strokes');
const darkModeMediaQuery = JSON.parse(document.querySelector('#darkModeMediaQuery').innerText);
const experimental       = JSON.parse(document.querySelector('#experimental').innerText);
const hole               = decodeURI(location.pathname.slice(1));
const langs              = JSON.parse(document.querySelector('#langs').innerText);
const picker             = document.querySelector('#picker');
const popups             = document.querySelector('#popups');
const restoreLink        = document.querySelector('#restoreLink');
const scorings           = ['Bytes', 'Chars'];
const scoringTabs        = document.querySelectorAll('#scoringTabs a');
const solutionPicker     = document.querySelector('#solutionPicker');
const solutions          = JSON.parse(document.querySelector('#solutions').innerText);
const status             = document.querySelector('#status');
const table              = document.querySelector('#scores');
const diff               = document.querySelector('#diff');
const sortedLangs        =
    Object.values(langs).sort((a, b) => a.name.localeCompare(b.name));

const darkMode = matchMedia(darkModeMediaQuery).matches;
let lang;
let latestSubmissionID = 0;
let solution = Math.max(scorings.indexOf(localStorage.getItem('solution')), 0);
let scoring  = Math.max(scorings.indexOf(localStorage.getItem('scoring')),  0);

// The savedInDB state is used to avoid saving solutions in localStorage when
// those solutions match the solutions in the database. It's used to avoid
// restoring a solution from localStorage when the user has improved that
// solution on a different browser. Assume the user is logged-in by default
// for non-experimental holes. At this point, it doesn't matter whether the
// user is actually logged-in, because solutions dictionaries will be empty
// for users who aren't logged-in, so the savedInDB state won't be used.
// By the time they are non-empty, the savedInDB state will have been updated.
let savedInDB = !experimental;

const keymap = JSON.parse(document.querySelector('#keymap').innerText);
const cm     = new CodeMirror(document.querySelector('#editor'), {
    autofocus:   true,
    indentUnit:  1,
    lineNumbers: true,
    smartIndent: false,
    theme:       darkMode ? 'material-ocean' : 'default',
    vimMode:     keymap == 'vim',
});

cm.on('change', () => {
    const code = cm.getValue();
    let infoText = '';
    for (let i = 0; i < scorings.length; i++) {
        if (i)
            infoText += ', ';
        infoText += `${comma(getScoring(code, i))} ${scorings[i].toLowerCase()}`;
    }
    chars.innerText = infoText;

    // Avoid future conflicts by only storing code locally that's different
    // from the server's copy.
    const serverCode = getSolutionCode(lang, solution);

    const key = getAutoSaveKey(lang, solution);
    if (code && (code !== serverCode || !savedInDB) && code !== langs[lang].example)
        localStorage.setItem(key, code);
    else
        localStorage.removeItem(key);

    updateRestoreLinkVisibility();
});

// Set/clear the hide-details cookie on details toggling.
document.querySelector('#details').ontoggle = e => document.cookie =
    'hide-details=;SameSite=Lax;Secure' + (e.target.open ? ';Max-Age=0' : '');

restoreLink.onclick = e => {
    cm.setValue(getSolutionCode(lang, solution));
    e.preventDefault();
};

(onhashchange = () => {
    lang = location.hash.slice(1) || localStorage.getItem('lang');

    // Kick 'em to Python if we don't know the chosen language.
    if (!langs[lang])
        lang = 'python';

    // Assembly only has bytes.
    if (lang == 'assembly')
        setSolution(0);

    localStorage.setItem('lang', lang);

    history.replaceState(null, '', '#' + lang);

    setCodeForLangAndSolution();
})();

// Wire up submit to clicking, keyboard, and maybe vim shortcut.
document.querySelector('#run a').onclick = submit;

onkeydown = e => (e.ctrlKey || e.metaKey) && e.key == 'Enter' ? submit() : undefined;

// Allow vim users to run code with :w or :write
if (cm.getOption('vimMode')) CodeMirror.Vim.defineEx('write', 'w', submit);

function getAutoSaveKey(lang, solution) {
    return `code_${hole}_${lang}_${solution}`;
}

function getOtherScoring(value) {
    return 1 - value;
}

function getScoring(str, index) {
    return scorings[index] == 'Bytes' ? byteLen(str) : charLen(str);
}

function getSolutionCode(lang, solution) {
    return lang in solutions[solution] ? solutions[solution][lang] : '';
}

function populateSolutionPicker() {
    while (solutionPicker.firstChild)
        solutionPicker.removeChild(solutionPicker.firstChild);

    const code0 = getSolutionCode(lang, 0);
    const code1 = getSolutionCode(lang, 1);
    const autoSave0 = localStorage.getItem(getAutoSaveKey(lang, 0));
    const autoSave1 = localStorage.getItem(getAutoSaveKey(lang, 1));

    // Only show the solution picker when both solutions are actually used.
    if (code0 && code1 && code0 != code1 || autoSave0 && autoSave1 && autoSave0 != autoSave1 ||
        // If a logged-in user has an auto-saved solution for the other
        // metric, that they have not submitted since logging in, they must be
        // allowed to switch to it, so they can submit it.
        (solution == 0 && code0 && autoSave1 && code0 != autoSave1) ||
        (solution == 1 && autoSave0 && code1 && autoSave0 != code1)) {
        for (let i = 0; i < scorings.length; i++) {
            let name = `Fewest ${scorings[i]}`;

            const solutionCode = getSolutionCode(lang, i);
            if (solutionCode) {
                name += ` <sup>${comma(getScoring(solutionCode, i))}</sup>`;
            }

            const child = document.createElement('a');
            child.id = `${scorings[i]}Solution`;
            child.innerHTML = name;
            if (i != solution) {
                child.href = '';
                child.onclick = e => {
                    e.preventDefault();
                    setSolution(i);
                    setCodeForLangAndSolution();
                };
            }
            solutionPicker.appendChild(child);
        }
    }
}

async function refreshScores() {
    // Rebuild the language picker with accurate stroke counts.
    picker.replaceChildren(...sortedLangs.map(l => {
        const tab = <a href={l.id == lang ? null : '#'+l.id}>{l.name}</a>;

        if (getSolutionCode(l.id, 0)) {
            const bytes = byteLen(getSolutionCode(l.id, 0));
            const chars = charLen(getSolutionCode(l.id, 1));

            let text = comma(bytes);
            if (chars && bytes != chars) text += '/' + comma(chars);

            tab.append(' ', <sup>{text}</sup>);
        }

        return tab;
    }));

    populateSolutionPicker();

    const scoringID = scorings[scoring].toLowerCase();
    const path      = `/${hole}/${lang}/${scoringID}`;
    const res       = await fetch(`/scores${path}/mini`);
    const rows      = res.ok ? await res.json() : [];

    document.querySelector('#allLink').href = '/rankings/holes' + path;

    table.replaceChildren(<tbody class={scoringID}>{
        // Rows.
        rows.map(r => <tr class={r.me ? 'me' : ''}>
            <td>{r.rank}<sup>{ord(r.rank)}</sup></td>
            <td>
                <a href={`/golfers/${r.login}`}>
                    <img src={`//avatars.githubusercontent.com/${r.login}?s=24`}/>
                    <span>{r.login}</span>
                </a>
            </td>
            <td data-tooltip={tooltip(r, 'Bytes')}>{comma(r.bytes)}</td>
            {r.chars !== null ?
                <td data-tooltip={tooltip(r, 'Chars')}>{comma(r.chars)}</td> : ''}
        </tr>)
    }{
        // Padding.
        [...Array(7 - rows.length).keys()].map(() =>
            <tr><td colspan="4">&nbsp;</td></tr>)
    }</tbody>);

    for (let i = 0; i < scoringTabs.length; i++) {
        const tab = scoringTabs[i];

        if (tab.innerText == scorings[scoring]) {
            tab.removeAttribute('href');
            tab.onclick = '';
        }
        else {
            tab.href = '';
            tab.onclick = e => {
                e.preventDefault();
                localStorage.setItem('scoring', scorings[scoring = i]);
                refreshScores();
            };
        }
    }
}

function setCodeForLangAndSolution() {
    if (solution != 0 && getSolutionCode(lang, 0) == getSolutionCode(lang, 1)) {
        const autoSave0 = localStorage.getItem(getAutoSaveKey(lang, 0));
        const autoSave1 = localStorage.getItem(getAutoSaveKey(lang, 1));
        if (autoSave0 && !autoSave1)
            setSolution(0);
    }

    const autoSaveCode = localStorage.getItem(getAutoSaveKey(lang, solution)) || '';
    const code = getSolutionCode(lang, solution) || autoSaveCode;

    cm.setOption('lineWrapping', lang != 'fish');
    cm.setOption('matchBrackets', lang != 'brainfuck' && lang != 'j' && lang != 'hexagony');
    cm.setOption('mode', {
        name: 'text/x-' + lang,
        startOpen: true,
        multiLineStrings: lang == 'c', // TCC supports multi-line strings
    });

    cm.setValue(autoSaveCode || code || langs[lang].example);

    refreshScores();

    for (const info of document.querySelectorAll('main .info'))
        info.style.display = info.classList.contains(lang) ? 'block' : '';
}

function setSolution(value) {
    localStorage.setItem('solution', scorings[solution = value]);
}

async function submit() {
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
    savedInDB = data.LoggedIn && !experimental;

    if (submissionID != latestSubmissionID)
        return;

    if (data.Pass) {
        for (let i = 0; i < scorings.length; i++) {
            const solutionCode = getSolutionCode(codeLang, i);
            if (!solutionCode || getScoring(code, i) <= getScoring(solutionCode, i)) {
                solutions[i][codeLang] = code;

                // Don't need to keep solution in local storage because it's
                // stored on the site. This prevents conflicts when the
                // solution is improved on another browser.
                if (savedInDB && localStorage.getItem(getAutoSaveKey(codeLang, i)) == code)
                    localStorage.removeItem(getAutoSaveKey(codeLang, i));
            }
        }
    }

    for (let i = 0; i < scorings.length; i++) {
        const key = getAutoSaveKey(codeLang, i);
        if (savedInDB) {
            // If the auto-saved code matches either solution, remove it to
            // avoid prompting the user to restore it.
            const autoSaveCode = localStorage.getItem(key);
            for (let j = 0; j < scorings.length; j++) {
                if (getSolutionCode(codeLang, j) == autoSaveCode)
                    localStorage.removeItem(key);
            }
        }
        else if (getSolutionCode(codeLang, i)) {
            // Autosave the best solution for each scoring metric, but don't
            // save two copies of the same solution, because that can lead to
            // the solution picker being show unnecessarily.
            if ((i == 0 || getSolutionCode(codeLang, 0) != getSolutionCode(codeLang, i)) &&
                getSolutionCode(codeLang, i) !== langs[codeLang].example)
                localStorage.setItem(key, getSolutionCode(codeLang, i));
            else
                localStorage.removeItem(key);
        }
    }

    // Automatically switch to the solution whose code matches the current
    // code after a new solution is submitted. Don't change scoring,
    // refreshScores will update the solution picker.
    if (data.Pass && getSolutionCode(codeLang, solution) != code &&
        getSolutionCode(codeLang, getOtherScoring(solution)) == code)
        setSolution(getOtherScoring(solution));

    // Update the restore link visibility, after possibly changing the active
    // solution.
    updateRestoreLinkVisibility();

    document.querySelector('h2').innerText = data.Pass ? 'Pass 😀' : 'Fail ☹️';

    // Show args if we have 'em.
    if (data.Argv.length) {
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

    const diffContent = document.querySelector('#diff-content');
    const diffVisible = attachDiff(diffContent, hole, data.Exp, data.Out, data.Argv, false);
    diff.style.display = diffVisible ? 'block' : 'none';

    status.className = data.Pass ? 'green' : 'red';
    status.style.display = 'block';

    refreshScores();

    // Show cheevos.
    popups.replaceChildren(...data.Cheevos.map(c => <div>
        <h3>Achievement Earned!</h3>
        { c.emoji }<p>{ c.name }</p>
    </div>));
}

function tooltip(row, scoring) {
    const bytes = scoring === 'Bytes' ? row.bytes : row.chars_bytes;
    const chars = scoring === 'Chars' ? row.chars : row.bytes_chars;

    return `${scoring} solution is ${comma(bytes)} bytes` +
        (chars !== null ? `, ${comma(chars)} chars.` : '.');
}

function updateRestoreLinkVisibility() {
    const serverCode = getSolutionCode(lang, solution);
    restoreLink.style.display =
        serverCode && cm.getValue() != serverCode ? 'initial' : 'none';
}
