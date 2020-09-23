const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

for (const a of document.querySelectorAll('.log-in')) {
    const url = new URL(a.href);

    // Add current time zome to the redirect URI.
    const redirect = new URL(url.searchParams.get('redirect_uri'));
    redirect.searchParams.set('time_zone', timeZone);
    url.searchParams.set('redirect_uri', redirect);

    a.href = url;
}
