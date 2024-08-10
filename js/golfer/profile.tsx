import { $ } from '../_util';

import { Chart, Legend, LineElement, PointElement, RadarController,
    RadialLinearScale, Tooltip } from 'chart.js';

Chart.register(Legend, LineElement, PointElement, RadarController,
    RadialLinearScale, Tooltip);

const style = getComputedStyle(document.body);

Chart.defaults.backgroundColor = style.getPropertyValue('--background');
Chart.defaults.borderColor     = style.getPropertyValue('--light-grey');
Chart.defaults.color           = style.getPropertyValue('--color');

const data = JSON.parse($('#category-overview-data').innerText);

new Chart($<HTMLCanvasElement>('canvas'), {
    type: 'radar',
    data: {
        labels: Object.keys(data.bytes),
        datasets: [{
            backgroundColor: style.getPropertyValue('--light-blue'),
            borderColor:     style.getPropertyValue('--blue'),
            data:            Object.values(data.bytes),
            label:           'Bytes',
        }, {
            backgroundColor: style.getPropertyValue('--light-green'),
            borderColor:     style.getPropertyValue('--green'),
            data:            Object.values(data.chars),
            label:           'Chars',
        }],
    },
    options: {
        scales: {
            r: {
                grid:  { color: style.getPropertyValue('--light-grey') },
                ticks: { showLabelBackdrop: false },
            },
        },
    },
});
