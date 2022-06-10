import { $ } from './_util';

$('#copyJson')?.addEventListener('click', () =>
    navigator.clipboard.writeText($('#dataJson').innerText));
$('#copyOutput')?.addEventListener('click', () =>
    navigator.clipboard.writeText($('#dataOutput').innerText));
