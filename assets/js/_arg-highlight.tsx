import { $$ }             from './_util';
import UnprintableElement from './_unprintable';

// Argument spans that, on hover, highlight their matching line in the
// expected and output panels (see argLines).
export function argSpans(argv: string[]) {
    return argv.map((arg, i) => {
        const span = <span>{arg}</span>;
        span.onmouseenter = () => highlight(i, true);
        span.onmouseleave = () => highlight(i, false);
        return span;
    });
}

// When there is one line per argument, wrap each line in a span tagged with
// its argument index. Exclude holes with one big argument like QR Decoder.
export function argLines(text: string, argv: string[], escape = false) {
    const render = (s: string) =>
        escape ? UnprintableElement.escape(s) : document.createTextNode(s);

    const lines = text.split('\n');
    if (argv.length < 2 || lines.length != argv.length) return render(text);

    const frag = document.createDocumentFragment();
    lines.forEach((line, i) => {
        if (i) frag.append('\n');
        frag.append(<span data-arg={String(i)}>{render(line)}</span>);
    });
    return frag;
}

function highlight(i: number, on: boolean) {
    for (const line of $$(`#exp [data-arg="${i}"], #out [data-arg="${i}"]`)) {
        line.classList.toggle('arg-highlight', on);
        if (on) center(line);
    }
}

// Scroll the line's panel, not the page, to keep the line in the middle.
// The panel is #exp/#out in the tabbed layout, or the div inside it otherwise.
function center(line: HTMLElement) {
    const panel = line.closest<HTMLElement>('#exp > div, #out > div, #exp, #out')!;
    const lineBox = line.getBoundingClientRect(), panelBox = panel.getBoundingClientRect();

    panel.scrollBy({
        behavior: 'smooth',
        top: lineBox.top - panelBox.top - (panel.clientHeight - lineBox.height) / 2,
    });
}
