'use strict';

/* include chart.js   */
/* include countUp.js */

Chart.defaults.global.layout.padding = 10;
Chart.defaults.global.legend.display = false;

onload = function() {
    for (let span of document.querySelectorAll("span"))
        new CountUp(span, 0, span.innerText.replace(/,/, '')).start();

    let canvases = document.querySelectorAll("canvas");

    for (let canvas of canvases) {
        canvas.height = canvas.parentNode.offsetHeight;
        canvas.width  = canvas.parentNode.offsetWidth;
    }

    new Chart(canvases[0].getContext('2d'), {
        type: 'pie',
        data: {
            labels: ['Fast', 'Medium', 'Slow'],
            datasets: [{
                data: JSON.parse(canvases[0].dataset.data),
                backgroundColor: ['#3d8b3d', '#df8a13', '#b52b27'],
            }],
        },
        options: { title: { display: true, text: 'Holes by Difficulty' } },
    });

    new Chart(canvases[1].getContext('2d'), {
        type: 'pie',
        data: {
            labels: ['Bash', 'JavaScript', 'Lua', 'Perl', 'Perl 6', 'PHP', 'Python', 'Ruby'],
            datasets: [{
                data: JSON.parse(canvases[1].dataset.data),
                backgroundColor: [
                    '#4dc9f6',
                    '#f67019',
                    '#f53794',
                    '#537bc4',
                    '#acc236',
                    '#166a8f',
                    '#00a950',
                    '#58595b',
                    '#8549ba',
                ],
            }],
        },
        options: { title: { display: true, text: 'Solutions by Language' } },
    });

    new Chart(canvases[2].getContext('2d'), {
        type: 'line',
        data: {
            datasets: [{
                data: JSON.parse(canvases[2].dataset.data),
                borderColor: "rgb(54, 162, 235)",
                backgroundColor: "rgb(54, 162, 235, .5)",
                pointBorderColor: "rgb(54, 162, 235)",
                pointBackgroundColor: "rgb(54, 162, 235)",
                pointBorderWidth: 1
            }]
        },
        options: {
            title: { display: true, text: 'Distribution of Holes Passed' },
            scales: { xAxes: [{ ticks: { min: 1 }, type: 'linear' }] },
            tooltips: {
                callbacks: {
                    label: function(item) {
                        let { xLabel: x, yLabel: y } = item;

                        // HACK https://stackoverflow.com/questions/46856815
                        if (x == 1.05)
                            x = 1;

                        return `${y} golfer${y == 1 ? ' has' : 's have'} passed ${x} hole${x == 1 ? '' : 's'}`;
                    },
                    title: function() {},
                },
                displayColors: false,
            }
        }
    });
};
