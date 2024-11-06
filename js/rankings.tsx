import { $ } from './_util';

import { Chart, BarController,
    LinearScale, PointElement, Tooltip, CategoryScale,
    BarElement } from 'chart.js';

const dataElement = $('#chart-data');
if (dataElement) {
    const scoring = location.pathname.split('/').filter(Boolean).pop();
    const {golfer, golfers} = JSON.parse(dataElement.innerText) as {golfer: number, golfers: {strokes: number, frequency: number}[]};
    if (golfers.length > 0) {
        const strokes = golfers.map(({strokes}) => strokes);
        const frequencies = golfers.map(({frequency}) => frequency);
        let total = 0;
        const cumsum = [0];
        for (let i = golfers.length - 1; i >= 0; i--) {
            cumsum.push(total += golfers[i].frequency);
        }
        cumsum.reverse();

        const orderedFrequencies = [...frequencies];
        orderedFrequencies.sort((a, b) => b - a);

        const golferBar = new Array(golfers.length);
        if (golfer) {
            golferBar[strokes.findIndex(x => x == golfer)] = 1 << 24;
        }

        const yMax = orderedFrequencies.length > 1 ? Math.min(orderedFrequencies[0], orderedFrequencies[1] * 1.5) : orderedFrequencies[0];

        Chart.register(BarController, BarElement,
            LinearScale, PointElement, Tooltip, CategoryScale);

        const style = getComputedStyle(document.body);

        Chart.defaults.backgroundColor = style.getPropertyValue('--background');
        Chart.defaults.borderColor     = style.getPropertyValue('--light-grey');
        Chart.defaults.color           = style.getPropertyValue('--color');

        const chart = new Chart($<HTMLCanvasElement>('#chart-container canvas'), {
            type: 'bar',
            data: {
                labels: strokes,
                datasets: [{
                    backgroundColor: style.getPropertyValue('--blue'),
                    barPercentage: 1,
                    categoryPercentage: 1,
                    data: frequencies,
                }, {
                    backgroundColor: style.getPropertyValue('--light-green'),
                    barPercentage: 1,
                    categoryPercentage: 1,
                    data: golferBar,
                }],
            },
            options: {
                responsive: true,
                aspectRatio: 6,
                plugins: {
                    legend: {display: false},
                    tooltip: {
                        callbacks: {
                            label(item) {
                                if (total < 2) return '';
                                const betterThan = cumsum[1 + item.dataIndex];
                                const betterThanRatio = betterThan / (total - 1);
                                return `Better than ${(betterThanRatio * 100).toFixed(2)}% golfers`;
                            },
                            title(items) {
                                const strokes = Number(items[0].label);
                                let otherGolfers = Number(items[0].raw);
                                if (golfer == strokes) otherGolfers--;
                                return `${strokes} ${scoring}: ${golfer == strokes ? 'you and ' : ''}${otherGolfers} golfer${otherGolfers > 1 ? 's' : ''}`;
                            },
                        },
                        filter(item) {
                            return item.datasetIndex < 1;
                        },
                        displayColors: false,
                        intersect: false,
                        mode: 'index',
                    },
                },
                scales: {
                    x: {
                        type: 'category',
                        display: false,
                        stacked: true,
                    },
                    y: {
                        type: 'linear',
                        max: yMax,
                        display: false,
                        beginAtZero: true,
                    },
                },
            },
        });
        chart.resize();
    }
}
