import { $, comma } from '../_util';

const form   = document.forms[0];
const run    = $('#run');
const status = $('#status');
const stop   = $('#stop');
const table  = $('table');
const tbody  = $('tbody');

const reduce = (acc: any, cur: any) => ({ ...acc, [cur.value]: cur.label });
const holes  = [...form.hole.options].reduce(reduce, {});
const langs  = [...(form.lang as any).options].reduce(reduce, {});

stop.onclick = () => alert('TODO');

form.onchange = () => history.replaceState(
    '', '', 'solutions?' + new URLSearchParams(new FormData(form) as any));

form.onsubmit = async e => {
    e.preventDefault();

    const start = Date.now();
    tbody.innerHTML = '';

    const res = await fetch('solutions/run?' +
        new URLSearchParams(new FormData(form) as any));

    if (!res.ok || !res.body) return;

    run.style.display   = 'none';
    stop.style.display  = 'block';
    table.style.display = 'table';

    let buffer     = '';
    let failing   = 0;
    let solutions = 0;

    const decoder = new TextDecoder();
    const reader  = res.body.getReader();
    const append  = (lineString: string) => {
        const line = JSON.parse(lineString);

        if (!line.pass) failing++;

        status.innerText = comma(++solutions) + '/' + comma(line.total) +
            ` solutions (${failing} failing) in ` +
            new Date(Date.now() - start).toISOString().substr(14, 8).replace(/^00:/, '');

        tbody.append(<tr>
            <td>{holes[line.hole]}</td>
            <td>{langs[line.lang]}</td>
            <td>{`${line.golfer} (${line.golfer_id})`}</td>
            <td>{comma(Math.round(line.took / 1e6)) + 'ms'}</td>
            <td><span class={line.pass ? 'green' : 'red'}>
                {line.pass ? 'PASS' : 'FAIL'}
            </span></td>
            <td><code>{line.stderr}</code></td>
        </tr>);
    };

    reader.read().then(function process({ done, value }) {
        if (done) {
            if (buffer)
                append(buffer);

            run.style.display  = 'block';
            stop.style.display = 'none';

            return;
        }

        const lines = (buffer += decoder.decode(value)).split(/\n(?=.)/);
        buffer = lines.pop() ?? '';
        lines.forEach(append);

        reader.read().then(process);
    });
};
