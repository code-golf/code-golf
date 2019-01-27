'use strict';

let nowDate  = new Date().toDateString();
let nowMonth = new Date().getMonth();
let nowYear  = new Date().getFullYear();
let months   = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];

for (let time of document.querySelectorAll('time')) {
    let date = new Date(time.getAttribute('datetime'));

    time.setAttribute('title', date);

    // If it happened today, show the time, zero padded.
    if (date.toDateString() == nowDate)
        time.innerText =
            ('0'+date.getHours()  ).slice(-2) + ':' +
            ('0'+date.getMinutes()).slice(-2);
    // Else show the date if it's within 11 months, otherwise the year.
    else {
        let diff =
            (nowYear - date.getFullYear()) * 12 + nowMonth - date.getMonth();

        time.innerText = diff < 11 ? date.getDate() + ' ' + months[date.getMonth()] : date.getFullYear();
    }
}
