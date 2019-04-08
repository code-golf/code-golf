'use strict';

const selects = document.querySelectorAll('select');
const change =
    () => location = '/scores/' + selects[0].value + '/' + selects[1].value;

for (const select of selects)
    select.onchange = change;
