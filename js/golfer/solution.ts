import { $ } from '../_util';

$('#test-solution').onclick = async () => {
    const res = await fetch('', { method: 'POST' });

    if (res.status != 200) {
        alert('Error ' + res.status);
        return;
    }

    let i = 0, msg = '';
    for (const run of await res.json()) {
        msg += `Run ${++i}: ${run.pass ? 'PASS' : 'FAIL'}`;
        msg += ` (${Math.round(run.time / 10**6)}ms`;
        if (run.timeout) msg += '; Timeout';
        msg += ')\n';
    }
    alert(msg);
};
