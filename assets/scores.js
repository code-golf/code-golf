'use strict';

const thisYear = new Date().getFullYear();
const today    = new Date().toDateString();
const months   = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];

onload = function() {
    for (let time of document.querySelectorAll('time')) {
        let date = new Date(time.getAttribute('datetime'));

        time.setAttribute('title', date.toString());

        time.innerText = date.getFullYear() != thisYear
                       ? date.getFullYear()
                       : date.toDateString() == today
                       ? ('0'+date.getHours()).slice(-2) + ':' + ('0'+date.getMinutes()).slice(-2)
                       : date.getDate() + ' ' + months[date.getMonth()];
    }

    let selects = document.querySelectorAll('select');

    for (let select of selects)
        select.onchange = function() {
            location = '/scores/' + selects[0].value + '/' + selects[1].value;
        };
};
