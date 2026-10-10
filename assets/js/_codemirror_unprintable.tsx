// CodeMirror unprintable character extensions
import { EditorView, keymap, Panel, showPanel, highlightSpecialChars } from '@codemirror/view';
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
    decimal?: boolean;
}

// Alt + hex digits inserts a character by hex code point, Shift + Alt + digits by decimal.
const updateInsertCharState = StateEffect.define<string | boolean>();
const startDecimalInsertChar = StateEffect.define<null>();
export const insertCharState = StateField.define<InsertCharState>({
    create: () => ({}),
    update(value, tr) {
        if (tr.docChanged) value = {};
        for (const e of tr.effects) {
            if (e.is(startDecimalInsertChar)) value = { code: '', decimal: true };
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
            const { code, toggleMode, decimal } = update.state.field(insertCharState);
            const toggleModeHelp = toggleMode ? ' (press Alt again to insert)' : '';
            if (decimal) {
                dom.textContent = code
                    ? `Composing decimal ${code}...`
                    : 'Type decimal digits and release Shift+Alt to insert.';
            }
            else if (code) {
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

        if (event.altKey && event.shiftKey && !event.ctrlKey && !event.metaKey) {
            const { decimal } = view.state.field(insertCharState);
            // Use the physical key, as Shift changes the digit row's characters on most layouts.
            const digit = event.code.match(/^(?:Digit|Numpad)(\d)$/)?.[1];

            if (event.key == 'Alt' || event.key == 'Shift') {
                if (!decimal) view.dispatch({ effects: startDecimalInsertChar.of(null) });
                event.preventDefault();
                return;
            }
            if (decimal && digit) {
                view.dispatch({ effects: updateInsertCharState.of(digit) });
                event.preventDefault();
                return;
            }
        }

        if (event.key == 'Alt' && !hasOtherMods) {
            // This might be a toggle start, so we start buffering keys.
            view.dispatch({ effects: updateInsertCharState.of('')});
            event.preventDefault();
            return;
        }

        const { code, toggleMode, decimal } = view.state.field(insertCharState);
        if ((event.altKey || toggleMode) && !hasOtherMods && event.key.match(/^[0-9a-f]$/i)) {
            view.dispatch({ effects: updateInsertCharState.of(event.key.toUpperCase()) });
            event.preventDefault();
        }
        else if (code || toggleMode || decimal) {
            // Reset the state whenever possible.
            view.dispatch({ effects: updateInsertCharState.of(false) });
        }
    },

    keyup: (event, view) => {
        const { code, toggleMode, decimal } = view.state.field(insertCharState);

        // Insert as soon as either Shift or Alt is released.
        if (decimal && (event.key == 'Alt' || event.key == 'Shift')) {
            let codepoint;
            try {
                if (code) codepoint = String.fromCodePoint(parseInt(code, 10));
            }
            catch {}

            view.dispatch(
                { effects: updateInsertCharState.of(false) },
                codepoint ? view.state.replaceSelection(codepoint) : {},
            );
            event.preventDefault();
            return;
        }

        if (event.key != 'Alt') return;

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

export const showUnprintables = highlightSpecialChars({
    specialChars: UnprintableElement.PATTERN,
    render: (code) => <u-p c={String.fromCharCode(code)} />,
});
