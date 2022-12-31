import { $, $$ } from './_util';
import dialogPolyfill from 'dialog-polyfill';

const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

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

// Polyfill dialog support, mainly for iOS.
for (const dialog of $$<HTMLDialogElement>('dialog')) {
    dialogPolyfill.registerDialog(dialog);

    dialog.onclick = e => e.target == dialog ? dialog.close() : null;
}

// Wire up any dialog buttons.
for (const btn of $$<HTMLElement>('[data-dialog]'))
    btn.onclick = () => {
        const dialog = $<HTMLDialogElement>('#' + btn.dataset.dialog);

        // If the dialog contains a form then reset it first.
        dialog.querySelector('form')?.reset();
        dialog.showModal();
    };
