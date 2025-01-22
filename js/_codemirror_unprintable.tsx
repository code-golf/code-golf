// CodeMirror unprintable character extensions
import { Decoration, EditorView, keymap, MatchDecorator, Panel, ViewPlugin, WidgetType, showPanel } from '@codemirror/view';
import { EditorState, StateEffect, StateField } from '@codemirror/state';
import UnprintableElement from './_unprintable';

export const carriageReturn = [
    EditorState.lineSeparator.of('\n'), // Prevent CM from treating carriage return as newline
    keymap.of({
        key: 'Shift-Enter',
        run: ({ state, dispatch }: any) => {
            dispatch(state.replaceSelection('\r'));
            return true;
        },
    } as any),
    /* When all the newlines inserted in a transaction are preceded by a
    carriage return, remove the carriage returns. This fixes lines ending
    with a carriage return when copied and pasted on Windows. */
    EditorState.transactionFilter.of(transaction => {
        const changes: {from: number, to: number}[] = [];
        let allPrefixed = true;
        transaction.changes.iterChanges((fromA, toA, fromB, toB, inserted) => {
            if (!allPrefixed)
                return;
            const string = inserted.sliceString(0);
            for (let i = 0; (i = string.indexOf('\n', i)) >= 0; i++) {
                if (string[i - 1] != '\r') {
                    allPrefixed = false;
                    return;
                }
                changes.push({ from: fromB + i - 1, to: fromB + i });
            }
        });
        return allPrefixed ? [transaction, { changes, sequential: true }] : transaction;
    }),
];

interface InsertCharState {
    code?: string;
    toggleMode?: boolean;
}

const updateInsertCharState = StateEffect.define<string | boolean>();
export const insertCharState = StateField.define<InsertCharState>({
    create: () => ({}),
    update(value, tr) {
        if (tr.docChanged) value = {};
        for (const e of tr.effects) {
            if (!e.is(updateInsertCharState)) continue;
            if (e.value === true) {
                value = { ...value, toggleMode: true };
            }
            else if (e.value === false) {
                value = {};
            }
            else {
                value = { ...value, code: (value.code || '') + e.value };
            }
        }
        return value;
    },
    provide: f => showPanel.from(f, value => value.code !== undefined ? createPanel : null),
});

const createPanel = (): Panel => {
    const dom = document.createElement('div');
    return {
        dom,
        update(update) {
            const { code, toggleMode } = update.state.field(insertCharState);
            const toggleModeHelp = toggleMode ? ' (press Alt again to insert)' : '';
            if (code) {
                dom.textContent = 'Composing \\u' + code + '...' + toggleModeHelp;
            }
            else if (toggleMode) {
                dom.textContent = 'Type hexadecimal digits and press Alt again to insert.';
            }
            else {
                dom.textContent = '';
            }
        },
    };
};

export const insertChar = EditorView.domEventHandlers({
    keydown: (event, view) => {
        const hasOtherMods = event.ctrlKey || event.shiftKey || event.metaKey;

        if (event.key == 'Alt' && !hasOtherMods) {
            // This might be a toggle start, so we start buffering keys.
            view.dispatch({ effects: updateInsertCharState.of('')});
            event.preventDefault();
            return;
        }

        const { code, toggleMode } = view.state.field(insertCharState);
        if ((event.altKey || toggleMode) && !hasOtherMods && event.key.match(/^[0-9a-f]$/i)) {
            view.dispatch({ effects: updateInsertCharState.of(event.key.toUpperCase()) });
            event.preventDefault();
        }
        else if (code || toggleMode) {
            // Reset the state whenever possible.
            view.dispatch({ effects: updateInsertCharState.of(false) });
        }
    },

    keyup: (event, view) => {
        if (event.key != 'Alt') return;

        const { code, toggleMode } = view.state.field(insertCharState);
        if (code === '' && !toggleMode) {
            view.dispatch({ effects: updateInsertCharState.of(true) });
            return;
        }

        let codepoint;
        try {
            if (code) codepoint = String.fromCodePoint(parseInt(code, 16));
        }
        catch {}

        view.dispatch(
            { effects: updateInsertCharState.of(false) },
            codepoint ? view.state.replaceSelection(codepoint) : {},
        );
    },
});

class UnprintableWidget extends WidgetType {
    value;

    constructor(value: number) {
        super();
        this.value = value;
    }
    toDOM() {
        return <u-p c={String.fromCharCode(this.value)} />;
    }
}

const unprintableDecorator = new MatchDecorator({
    regexp: UnprintableElement.PATTERN,
    decoration: match => Decoration.replace({
        widget: new UnprintableWidget(match[0].charCodeAt(0)),
    }),
});

export const showUnprintables = ViewPlugin.fromClass(
    class {
        decorations;
        constructor(view: any) {
            this.decorations = unprintableDecorator.createDeco(view);
        }
        update(update: any) {
            this.decorations = unprintableDecorator.updateDeco(update, this.decorations);
        }
    },
    { decorations: plugin => plugin.decorations },
);
