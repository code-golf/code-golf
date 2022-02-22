document.querySelector('#copy')?.addEventListener('click', () =>
    navigator.clipboard.writeText(document.querySelector('#data').innerText));
