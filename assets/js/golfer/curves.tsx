import { $, $$, comma } from '../_util';

import { Chart, Legend, LinearScale, LineController, LineElement,
    LogarithmicScale, PointElement, TimeScale, Tooltip } from 'chart.js';

import 'chartjs-adapter-luxon';

Chart.register(Legend, LinearScale, LineController, LineElement,
    LogarithmicScale, PointElement, TimeScale, Tooltip);

interface Improvement {
    bytes: number, chars: number | null, hole: string, lang: string,
    scoring: 'bytes' | 'chars', submitted: string,
}

interface Estimator { hole: string, lang: string, scoring: string, sb: number }

interface Medal { lang: string, medal: string, scoring: string, submitted: string }

const data: {
    cheevos: string[] | null, estimators: Estimator[] | null,
    improvements: Improvement[] | null, langNames: Record<string, string>,
    medals: Medal[] | null,
} = JSON.parse($('#curves-data').innerText);

const style = getComputedStyle(document.body);
const color = (name: string) => style.getPropertyValue(name).trim();

Chart.defaults.borderColor = color('--light-grey');
Chart.defaults.color       = color('--color');

// Points share the colour of their scoring, dashed. Medals use their own
// colours, the rest are told apart by dash pattern.
const series = [
    { id: 'bytes',       label: 'Bytes',        color: '--blue-fg'   },
    { id: 'chars',       label: 'Chars',        color: '--green-fg'  },
    { id: 'bytesPoints', label: 'Bytes points', color: '--blue-fg',   dash: [6, 4] },
    { id: 'charsPoints', label: 'Chars points', color: '--green-fg',  dash: [6, 4] },
    { id: 'holes',       label: '⛳ Holes',     color: '--red-fg'    },
    { id: 'langs',       label: '🔣 Langs',     color: '--color',     dash: [2, 3] },
    { id: 'cheevos',     label: '🏆 Cheevos',   color: '--color',     dash: [8, 3, 2, 3] },
    { id: 'unicorn',     label: '🦄 Unicorns',  color: '--purple-fg' },
    { id: 'diamond',     label: '💎 Diamonds',  color: '--diamond'   },
    { id: 'gold',        label: '🥇 Gold',      color: '--yellow-fg' },
    { id: 'silver',      label: '🥈 Silver',    color: '--grey'      },
    { id: 'bronze',      label: '🥉 Bronze',    color: '--orange-fg' },
] as const;

type SeriesID = typeof series[number]['id'];
type Point    = { x: number, y: number | null };

const langSelect = $<HTMLSelectElement>('#lang');
langSelect.append(...Object.entries(data.langNames)
    .sort(([, a], [, b]) => a.localeCompare(b))
    .map(([id, name]) => <option value={id}>{name}</option>));

const checkbox = (id: SeriesID) => $<HTMLInputElement>(`#series [name=${id}]`);

// Restore the filters from the URL so a set of curves can be shared.
const params = new URLSearchParams(location.search);
if (params.has('lang') && params.get('lang')! in data.langNames)
    langSelect.value = params.get('lang')!;
if (params.has('curves'))
    for (const { id } of series)
        checkbox(id).checked = params.get('curves')!.split(',').includes(id);

const day = (date: string) => new Date(date).setHours(0, 0, 0, 0);

// Replay the golfer's history, one point per day for every series.
function computeSeries(lang: string) {
    const matches = (l: string) => !lang || l === lang;

    const sb = new Map((data.estimators ?? []).map(e =>
        [`${e.hole} ${e.lang} ${e.scoring}`, e.sb]));

    type Event = { time: number, apply: () => void };
    const events: Event[] = [];

    const totals = Object.fromEntries(series.map(({ id }) => [id, 0])) as Record<SeriesID, number>;

    const solutions  = new Map<string, Improvement>();
    const holeLangs  = new Map<string, Set<string>>();
    const holePoints = new Map<string, number>();
    const langsUsed  = new Set<string>();

    // Points on today's leaderboards, best lang per hole like the profile.
    const pointsFor = (hole: string, scoring: string) => {
        let best = 0;
        for (const solution of solutions.values()) {
            if (solution.hole !== hole || solution.scoring !== scoring) continue;

            const estimator = sb.get(`${hole} ${solution.lang} ${scoring}`) ?? 0;
            const strokes   = scoring === 'bytes' ? solution.bytes : solution.chars!;
            best = Math.max(best, Math.round(estimator / strokes * 1000));
        }
        return best;
    };

    for (const improvement of data.improvements ?? []) {
        if (!matches(improvement.lang)) continue;

        events.push({ time: day(improvement.submitted), apply: () => {
            const { hole, lang, scoring } = improvement;
            const key = `${hole} ${lang} ${scoring}`;
            const old = solutions.get(key);

            if (scoring === 'bytes')
                totals.bytes += improvement.bytes - (old?.bytes ?? 0);
            else if (improvement.chars !== null)
                totals.chars += improvement.chars - (old?.chars ?? 0);

            solutions.set(key, improvement);

            if (!holeLangs.has(hole)) holeLangs.set(hole, new Set);
            holeLangs.get(hole)!.add(lang);
            langsUsed.add(lang);
            totals.holes = holeLangs.size;
            totals.langs = langsUsed.size;

            const pointsID = scoring === 'bytes' ? 'bytesPoints' : 'charsPoints';
            const points   = pointsFor(hole, scoring);
            totals[pointsID] += points - (holePoints.get(`${hole} ${scoring}`) ?? 0);
            holePoints.set(`${hole} ${scoring}`, points);
        } });
    }

    // Medals of both scorings, like the profile.
    for (const medal of data.medals ?? [])
        if (matches(medal.lang))
            events.push({ time: day(medal.submitted),
                apply: () => totals[medal.medal as SeriesID]++ });

    // Cheevos aren't per lang.
    if (!lang)
        for (const earned of data.cheevos ?? [])
            events.push({ time: day(earned), apply: () => totals.cheevos++ });

    events.sort((a, b) => a.time - b.time);

    const points = Object.fromEntries(series.map(({ id }) =>
        [id, [] as Point[]])) as Record<SeriesID, Point[]>;

    // A log scale has no zero, so a curve starts at its first non-zero value.
    const snapshot = (time: number) => {
        for (const { id } of series)
            points[id].push({ x: time, y: totals[id] || null });
    };

    for (let i = 0; i < events.length; i++) {
        events[i].apply();

        // Only keep the last state of each day.
        if (events[i + 1]?.time !== events[i].time) snapshot(events[i].time);
    }

    // Carry the curves on until today.
    if (events.length) snapshot(day(new Date().toISOString()));

    return points;
}

const chart = new Chart<'line', Point[]>($<HTMLCanvasElement>('#chart canvas'), {
    type: 'line',
    data: {
        datasets: series.map(s => ({
            backgroundColor:  color(s.color),
            borderColor:      color(s.color),
            borderDash:       'dash' in s ? [...s.dash] : [],
            borderWidth:      2,
            data:             [],
            label:            s.label,
            pointHoverRadius: 4,
            pointRadius:      0,
            stepped:          true,
        })),
    },
    options: {
        animation:           false,
        interaction:         { intersect: false, mode: 'index' },
        maintainAspectRatio: false,
        plugins: {
            legend: {
                labels: { filter: item => !item.hidden },

                // The checkboxes below are the source of truth.
                onClick: (_e, item) => {
                    const { id } = series.find(s => s.label === item.text)!;
                    checkbox(id).checked = false;
                    update();
                },
            },
            tooltip: {
                callbacks: {
                    label: item => `${item.dataset.label}: ${comma(item.parsed.y ?? 0)}`,
                },
            },
        },
        scales: {
            x: {
                grid: { color: color('--light-grey') },
                time: { tooltipFormat: 'DD' },
                type: 'time',
            },
            y: {
                grid:  { color: color('--light-grey') },
                min:   1,
                ticks: {
                    // Only label the powers of ten, the rest is noise.
                    callback: value => Number.isInteger(Math.log10(+value)) ? comma(+value) : '',
                },
                type: 'logarithmic',
            },
        },
    },
});

function update() {
    const lang = langSelect.value;

    // Langs is always 1 when a single lang is picked, cheevos aren't per lang.
    checkbox('langs').disabled = checkbox('cheevos').disabled = !!lang;

    const points  = computeSeries(lang);
    const checked = series.filter(({ id }) => checkbox(id).checked).map(s => s.id);
    const shown   = checked.filter(id => !checkbox(id).disabled);

    series.forEach(({ id }, i) => {
        chart.data.datasets[i].data   = points[id];
        chart.data.datasets[i].hidden = !shown.includes(id);
    });
    chart.update();

    $('#chart').classList.toggle('hide', !shown.length);
    $('#no-curves').classList.toggle('hide', shown.length > 0);

    const url = new URL(location.href);
    url.search = new URLSearchParams({
        ...(lang ? { lang } : {}), curves: checked.join(','),
    }).toString();
    history.replaceState(null, '', url);
}

langSelect.onchange = update;
$$('#series input').forEach(input => input.onchange = update);

update();
