const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

// Add current time zone to the redirect URI of any log in links.
for (const a of document.querySelectorAll('.log-in')) {
    const url = new URL(a.href);

    const redirect = new URL(url.searchParams.get('redirect_uri'));
    redirect.searchParams.set('time_zone', timeZone);
    url.searchParams.set('redirect_uri', redirect);

    a.href = url;
}

// Wire up mobile form navigation.
document.querySelector('#form-nav')?.addEventListener('change', e =>
    location = [...new FormData(e.target.form).values()].filter(v =>
        v.length).join('/'));

// Add suggestions to any input with a list.
for (const input of document.querySelectorAll('[list]')) {
    let controller;

    input.oninput = async () => {
        controller?.abort();        // Abort the old request (if exists).
        input.list.innerHTML = '';  // Clear current suggestions.

        if (input.value != '')
            input.list.append(...(await (await fetch(
                `/api/suggestions/${input.list.id}?` +
                    new URLSearchParams({ ...input.dataset, q: input.value }),
                { signal: (controller = new AbortController()).signal },
            )).json()).map(suggestion => <option value={suggestion}/>));
    };
}
