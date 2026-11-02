// CodeMirror rainbow brackets: while Alt is held, colour (), [] and {} by nesting depth.
import { Decoration, DecorationSet, EditorView, ViewPlugin, ViewUpdate } from '@codemirror/view';
import { EditorState, RangeSetBuilder, StateEffect }                     from '@codemirror/state';
import { syntaxTree }                                                    from '@codemirror/language';

const pairs: Record<string, string> = { ')': '(', ']': '[', '}': '{' };
const depths = 5;
const marks  = [...Array(depths).keys()].map(i => Decoration.mark({ class: `cm-rainbow-${i + 1}` }));

const setShown = StateEffect.define<boolean>();

// Brackets inside strings, comments, etc. don't count.
const ignored = /string|comment|regexp|char/i;

function isIgnored(state: EditorState, pos: number) {
    for (let node: any = syntaxTree(state).resolveInner(pos, 1); node; node = node.parent)
        if (ignored.test(node.name)) return true;
    return false;
}

// Golf solutions are short, so scan the whole document rather than the viewport.
// Unmatched brackets are left uncoloured.
function buildDecorations(state: EditorState) {
    const depthAt: (number | undefined)[] = [];
    const stack: number[] = [];
    const doc = state.doc.toString();

    for (let pos = 0; pos < doc.length; pos++) {
        const char = doc[pos];
        const open = pairs[char];
        if (!open && !'([{'.includes(char) || isIgnored(state, pos)) continue;

        if (!open) stack.push(pos);
        else if (stack.length && doc[stack.at(-1)!] == open)
            depthAt[stack.pop()!] = depthAt[pos] = stack.length;
    }

    const builder = new RangeSetBuilder<Decoration>();
    depthAt.forEach((depth, pos) => builder.add(pos, pos + 1, marks[depth! % depths]));
    return builder.finish();
}

// Key events are only watched, never prevented, so the other Alt shortcuts keep working.
// AltGr (Ctrl + Alt on Windows) types brackets on many layouts, so it doesn't count.
export const rainbowBrackets = ViewPlugin.fromClass(class {
    decorations: DecorationSet = Decoration.none;
    shown = false;
    view: EditorView;

    sync = (event: Event) => {
        const shown = event instanceof KeyboardEvent && event.altKey
            && !event.ctrlKey && !event.getModifierState('AltGraph');
        if (shown != this.shown) this.view.dispatch({ effects: setShown.of(shown) });
    };

    constructor(view: EditorView) {
        this.view = view;
        addEventListener('keydown', this.sync, true);
        addEventListener('keyup',   this.sync, true);
        addEventListener('blur',    this.sync);
    }

    update(update: ViewUpdate) {
        let { shown } = this;
        for (const tr of update.transactions)
            for (const effect of tr.effects)
                if (effect.is(setShown)) shown = effect.value;

        if (shown != this.shown || shown && (update.docChanged
            || syntaxTree(update.startState) != syntaxTree(update.state)))
            this.decorations = shown ? buildDecorations(update.state) : Decoration.none;

        this.shown = shown;
    }

    destroy() {
        removeEventListener('keydown', this.sync, true);
        removeEventListener('keyup',   this.sync, true);
        removeEventListener('blur',    this.sync);
    }
}, { decorations: v => v.decorations });
