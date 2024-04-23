import { $ } from './_util';

import { Chart, LineController, LineElement,
    LinearScale, PointElement, TimeScale, Tooltip } from 'chart.js';

import 'chartjs-adapter-luxon';

const data = $('#chart-data');
if (data) {
    Chart.register(LineController, LineElement,
        LinearScale, PointElement, TimeScale, Tooltip);

    const style = getComputedStyle(document.body);

    Chart.defaults.backgroundColor = style.getPropertyValue('--background');
    Chart.defaults.borderColor     = style.getPropertyValue('--light-grey');
    Chart.defaults.color           = style.getPropertyValue('--color');

    const chart = new Chart($<HTMLCanvasElement>('#chart-container canvas'), {
        type: 'line',
        data: {
            datasets: [{
                backgroundColor: style.getPropertyValue('--light-blue'),
                borderColor:     style.getPropertyValue('--blue'),
                data:            JSON.parse(data.innerText),
            }],
        },
        options: {
            scales: {
                x: {
                    grid: { color: style.getPropertyValue('--light-grey') },
                    time: { unit: 'year' },
                    type: 'time',
                },
                y: {
                    grid: { color: style.getPropertyValue('--light-grey') },
                },
            },
        },
    });

    onresize = () => chart.resize();
}
