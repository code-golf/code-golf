import { $ } from './_inject';

$('#copy')?.addEventListener('click', () =>
    navigator.clipboard.writeText(($('#data') as HTMLElement).innerText));
