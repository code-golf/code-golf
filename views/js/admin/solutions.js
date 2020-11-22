const form   = document.forms[0];
const run    = document.querySelector('#run');
const status = document.querySelector('#status');
const stop   = document.querySelector('#stop');
const table  = document.querySelector('table');
const tbody  = document.querySelector('tbody');

const reduce = (acc, cur) => ({ ...acc, [cur.value]: cur.label });
const holes  = [...form.hole.options].reduce(reduce, {});
const langs  = [...form.lang.options].reduce(reduce, {});

const pass = document.createElement('span');
pass.className = 'green';
pass.innerText = 'PASS';

const fail = document.createElement('span');
fail.className = 'red';
fail.innerText = 'FAIL';

function make(tag, ...children) {
    const node = document.createElement(tag);
    node.append(...children);
    return node;
}

stop.onclick = () => alert('TODO');

form.onchange = () => history.replaceState(
    '', '', 'solutions?' + new URLSearchParams(new FormData(form)));

form.onsubmit = async e => {
    e.preventDefault();

    const start = Date.now();
    tbody.innerHTML = '';

    const res = await fetch('solutions/run?' +
        new URLSearchParams(new FormData(form)));

    if (!res.ok) return;

    run.style.display   = 'none';
    stop.style.display  = 'block';
    table.style.display = 'table';

    let buffer     = '';
    let failing   = 0;
    let solutions = 0;

    const decoder = new TextDecoder();
    const reader  = res.body.getReader();
    const append  = line => {
        line = JSON.parse(line);

        if (!line.pass) failing++;

        status.innerText = (++solutions).toLocaleString('en') +
            ` solutions (${failing} failing) in ` +
            new Date(Date.now() - start).toISOString().substr(14, 8).replace(/^00:/, '');

        tbody.append(
            make(
                'tr',
                make('td', holes[line.hole]),
                make('td', langs[line.lang]),
                make('td', `${line.golfer} (${line.golfer_id})`),
                make('td', Math.round(line.took / 1e6).toLocaleString('en') + 'ms'),
                make('td', (line.pass ? pass : fail).cloneNode(true)),
                make('td', make('code', line.stderr)),
            ),
        );
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
        buffer = lines.pop();
        lines.forEach(append);

        reader.read().then(process);
    });
};
