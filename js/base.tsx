import { requestResults, SearchNavResult } from './_search-nav';
import { $, $$ } from './_util';

// Work around a Chrome bug, force SVGs to re-appear, should be fixed in 144.
// See https://issues.chromium.org/issues/459746761 for details.
$$('use[href]').forEach(e => e.setAttribute('href', e.getAttribute('href')!));

const { timeZone } = Intl.DateTimeFormat().resolvedOptions();

// Add current time zone to the redirect URI of any log in links.
for (const a of $$<HTMLAnchorElement>('.log-in')) {
    const url = new URL(a.href);

    // Assume a redirect is already present
    const redirect = new URL(url.searchParams.get('redirect_uri') as string);
    redirect.searchParams.set('time_zone', timeZone);
    url.searchParams.set('redirect_uri', redirect.href);

    a.href = url.href;
}

// Wire up mobile form navigation.
$('#form-nav')?.addEventListener('change',
    (e: Event) => location.href =
        [
            ...new FormData((e.target as HTMLFormElement).form).values(),
        ].filter((v: string | FormDataEntryValue) => v).join('/'));

// Add suggestions to any input with a list.
for (const input of $$<any>('[list]')) {
    let controller: AbortController | undefined;

    input.oninput = async () => {
        controller?.abort();        // Abort the old request (if exists).
        input.list.innerHTML = '';  // Clear current suggestions.

        if (input.value != '')
            input.list.append(...(await (await fetch(
                `/api/suggestions/${input.list.id}?` +
                new URLSearchParams({ ...input.dataset, q: input.value }),
                { signal: (controller = new AbortController()).signal },
            )).json()).map((suggestion: string) => <option value={suggestion} />));
    };
}

// Safari 26.2 doesn't support but has it on the proto, so test attr instead.
// See https://github.com/tak-dcxi/dialog-closedby-polyfill/issues/13
const closedByAnySupported = (() => {
    try {
        const dialog = <dialog closedby="any"></dialog>;
        return dialog.closedBy === 'any';
    }
    catch {
        return false;
    }
})();

// Polyfill closedby="any"
// https://caniuse.com/mdn-html_elements_dialog_closedby
// onmousedown not onclick, see https://stackoverflow.com/questions/25864259
if (!closedByAnySupported)
    for (const dialog of $$<HTMLDialogElement>('dialog[closedby="any"]'))
        dialog.onmousedown = e => e.target == dialog ? dialog.close() : null;

// Polyfill command="show-modal"
// https://caniuse.com/mdn-api_commandevent_command
if (!('command' in HTMLButtonElement.prototype))
    for (const btn of $$('[command="show-modal"]'))
        btn.onclick = () =>
            $<HTMLDialogElement>('#' + btn.getAttribute('commandfor')).showModal();

// Reset forms inside dialogs on dialog close.
for (const dialog of $$('dialog:has(form)'))
    dialog.onclose = () => dialog.querySelector('form')!.reset();

// Search navigation dialog
document.addEventListener('keydown', e => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
        const dialog = $<HTMLDialogElement>('#search-nav-dialog');
        if (!dialog) return;

        if (dialog.open) {
            dialog.close();
        }
        else {
            dialog.showModal();
            updateResults([]);
            e.preventDefault();
        }
    }
});

window.addEventListener('hashchange', () => {
    const dialog = $<HTMLDialogElement>('#search-nav-dialog');
    dialog?.close();
});

$('#search-nav-input').onkeyup = () => requestResults($<HTMLInputElement>('#search-nav-input').value, updateResults);

$('#search-nav-form').onsubmit = e => {
    $<HTMLAnchorElement>('li.current-result a')?.click();
    e.preventDefault();
};

let currentIndex = 0;
let currentResults: SearchNavResult[] = [];

$<HTMLDialogElement>('#search-nav-dialog').onkeydown = (e: KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
        currentIndex++;
        renderResults();
        e.preventDefault();
    }
    else if (e.key === 'ArrowUp') {
        currentIndex--;
        renderResults();
        e.preventDefault();
    }
};

function updateResults(results: SearchNavResult[]) {
    currentResults = results;
    currentIndex = 0;
    renderResults();
}

const mod = (x: number, modulus: number) => (x % modulus + modulus) % modulus;

function renderResults() {
    const resultNodes = currentResults.map((r, i) => (<li class={i === mod(currentIndex, currentResults.length) ? 'current-result' : ''}><a href={r.path}><span class='result-description'>{r.description}</span> <span class='result-path'>{r.path}</span></a></li>));
    $('#search-nav-results').replaceChildren(...resultNodes);
}
