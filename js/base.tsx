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
            )).json()).map((suggestion: string) => <option value={suggestion}/>));
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

document.addEventListener('keydown', (e) => {
    if (e.ctrlKey && e.key === 'p') {
        const dialog = $<HTMLDialogElement>('#search-nav-dialog');
        if (!dialog) return;

        if (dialog.open) {
            dialog.close();
        }
        else {
            dialog.showModal();
            e.preventDefault();
        }
    }
});

$('#search-nav-input').onkeyup = requestResults;

interface SearchNavResult {
    path: string;
    description: string;
}
function updateResults(results: SearchNavResult[]) {
    const resultNodes = results.map(r => (<li><a href={r.path}>{r.description}</a></li>));
    $('#search-nav-results').replaceChildren(...resultNodes);
}

function requestResults() {
    updateResults([
        {path: "/day-of-week#nim", description: "Day of Week in Nim"},
        {path: "/golfers/MichalMarsalek", description: "My profile"},
    ])
}