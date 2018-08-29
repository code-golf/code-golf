'use strict';

/* include chart.js */

Chart.defaults.global.animation.duration = 2000;

Chart.defaults.global.layout.padding = 10;

Chart.defaults.global.legend.labels.fontColor =
Chart.defaults.global.title.fontColor = '#000';

Chart.defaults.global.legend.labels.fontFamily =
Chart.defaults.global.title.fontFamily = '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif';

onload = function() {
    for (let span of document.querySelectorAll('span')) {
        let count, end = Number(span.dataset.x), start;

        count = function(timestamp) {
            if (!start)
                start = timestamp;

            let progress = timestamp - start;

            let val = Math.floor(end * (progress / 1500));

            span.textContent =
                Math.min(val, end).toString().replace(/(\d)(\d{3})$/, '$1,$2');

            if (progress < 1500)
                requestAnimationFrame(count);
        };

        requestAnimationFrame(count);
    }

    let canvases = document.querySelectorAll('canvas');

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
                backgroundColor: ['#2ECC71', '#FFEB3B', '#E74C3C'],
            }],
        },
        options: {
            legend: { position: 'right' },
            title: { display: true, text: 'Holes by Difficulty' },
        },
    });

    let data = JSON.parse(canvases[1].dataset.data);

    new Chart(canvases[1].getContext('2d'), {
        type: 'pie',
        data: {
            // FIXME Maybe the DB enums should be in the correct case?
            labels: data[0].map(
                lang => lang[0].toUpperCase() + lang.slice(1).replace('hp', 'HP').replace('sc', 'Sc').replace(6, ' 6'),
            ),
            datasets: [{
                data: data[1],
                backgroundColor: [
                    '#3498DB',
                    '#E74C3C',
                    '#F39C12',
                    '#2ECC71',
                    '#8E44AD',
                    '#FFEB3B',
                    '#16A085',
                    '#7F8C8D',
                ],
            }],
        },
        options: {
            legend: { position: 'right' },
            title: { display: true, text: 'Solutions by Language' },
        },
    });

    new Chart(canvases[2].getContext('2d'), {
        type: 'line',
        data: {
            datasets: [{
                data: JSON.parse(canvases[2].dataset.data),
                borderColor: '#3498DB',
                backgroundColor: 'rgb(52, 152, 219, .5)',
                pointBorderColor: '#3498DB',
                pointBackgroundColor: '#3498DB',
                pointBorderWidth: 1,
            }]
        },
        options: {
            legend: { display: false },
            scales: { xAxes: [{ ticks: { min: 1 }, type: 'linear' }] },
            title: { display: true, text: 'Distribution of Holes Passed' },
            tooltips: {
                callbacks: {
                    label: function(item) {
                        let { xLabel: x, yLabel: y } = item;

                        return `${y} golfer${y == 1 ? ' has' : 's have'} passed ${x} hole${x == 1 ? '' : 's'}`;
                    },
                    title: function() {},
                },
                displayColors: false,
            }
        }
    });

    new Chart(canvases[3].getContext('2d'), {
        type: 'line',
        data: {
            datasets: [{
                data: JSON.parse(canvases[3].dataset.data),
                borderColor: '#3498DB',
                backgroundColor: 'rgb(52, 152, 219, .5)',
                pointBorderColor: '#3498DB',
                pointBackgroundColor: '#3498DB',
                pointBorderWidth: 1,
            }]
        },
        options: {
            legend: { display: false },
            scales: { xAxes: [{ ticks: { min: 1 }, type: 'linear' }] },
            title: { display: true, text: 'Distribution of Languages Used' },
            tooltips: {
                callbacks: {
                    label: function(item) {
                        let { xLabel: x, yLabel: y } = item;

                        return `${y} golfer${y == 1 ? ' has' : 's have'} used ${x} language${x == 1 ? '' : 's'}`;
                    },
                    title: function() {},
                },
                displayColors: false,
            }
        }
    });
};
