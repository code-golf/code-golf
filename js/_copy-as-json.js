$('#copy')?.addEventListener('click', () =>
    navigator.clipboard.writeText($('#data').innerText));
