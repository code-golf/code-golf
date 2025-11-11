import { requestResults, SearchNavResult } from './_search-nav';
import { $, $$ } from './_util';

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

// Close dialogs when clicking outside of them.
// onmousedown not onclick, see https://stackoverflow.com/questions/25864259
for (const dialog of $$<HTMLDialogElement>('dialog'))
    dialog.onmousedown = e => e.target == dialog ? dialog.close() : null;

// Wire up any dialog buttons.
for (const btn of $$<HTMLElement>('[data-dialog]'))
    btn.onclick = () => {
        const dialog = $<HTMLDialogElement>('#' + btn.dataset.dialog);

        // If the dialog contains a form then reset it first.
        dialog.querySelector('form')?.reset();
        dialog.showModal();
    };

// Search navigation dialog
document.addEventListener('keydown', e => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
        const dialog = $<HTMLDialogElement>('#search-nav-dialog');
        if (!dialog) return;

        if (dialog.open) {
            dialog.close();
        }
        else {
            dialog.querySelector('form')?.reset();
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