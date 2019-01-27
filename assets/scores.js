'use strict';

let checkbox = document.querySelector('input');
let selects = document.querySelectorAll('select');

for (let select of selects)
    select.onchange = function() {
        let url = '/scores';

        if (selects[0].value)
            url += '/' + selects[0].value;

        if (selects[1].value)
            url += '/' + selects[1].value;

        if (checkbox && checkbox.checked)
            url += '/all';

        location = url;
    };

if (checkbox)
    checkbox.onchange = selects[0].onchange;
