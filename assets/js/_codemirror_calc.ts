import { EditorSelection, StateEffect, StateField } from '@codemirror/state';
import { Decoration, EditorView }                   from '@codemirror/view';
import { calc }                                     from './_calc';

type Range = { from: number, to: number };

// Selections that failed to evaluate flash red briefly.
const setErrors = StateEffect.define<Range[]>();
const errorMark = Decoration.mark({ class: 'cm-calc-error' });
const timers    = new WeakMap<EditorView, ReturnType<typeof setTimeout>>();

const errors = StateField.define({
    create: () => Decoration.none,
    update(deco, tr) {
        for (const e of tr.effects)
            if (e.is(setErrors)) return Decoration.set(e.value.map(r => errorMark.range(r.from, r.to)));
        return deco.map(tr.changes);
    },
    provide: f => EditorView.decorations.from(f),
});

// Replace every selected expression with its value, in a single undo step.
export const evaluateSelections = (view: EditorView) => {
    const { state } = view, failed: Range[] = [];
    if (state.selection.ranges.every(r => r.empty)) return false;

    const spec = state.changeByRange(range => {
        if (range.empty) return { range };
        try {
            const insert = calc(state.sliceDoc(range.from, range.to));
            return {
                changes: { from: range.from, to: range.to, insert },
                range:   EditorSelection.range(range.from, range.from + insert.length),
            };
        }
        catch {
            failed.push(range);
            return { range };
        }
    });

    const { changes } = spec;
    view.dispatch({
        ...spec, scrollIntoView: true, userEvent: 'input.calc',
        effects: setErrors.of(failed.map(r => ({ from: changes.mapPos(r.from, 1), to: changes.mapPos(r.to, -1) }))),
    });

    if (failed.length) {
        clearTimeout(timers.get(view));
        timers.set(view, setTimeout(() => view.dom.isConnected &&
            view.dispatch({ effects: setErrors.of([]) }), 400));
    }

    return true;
};

export const calcErrors = [
    errors,
    EditorView.baseTheme({ '.cm-calc-error': { background: 'rgba(255, 0, 0, .35)' } }),
];
