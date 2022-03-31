import { $ } from './_inject';

$('#copy')?.addEventListener('click', () =>
    navigator.clipboard.writeText($('#data').innerText));
