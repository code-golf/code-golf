import { $ } from './_util';

$('#copy')?.addEventListener('click', () =>
    navigator.clipboard.writeText($('#data').innerText));
