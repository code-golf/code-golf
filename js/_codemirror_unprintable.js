// CodeMirror unprintable character extensions
import { Decoration, EditorView, keymap, MatchDecorator, ViewPlugin, WidgetType } from "@codemirror/view";
import { EditorState }                                                            from "@codemirror/state";

export const carriageReturn = [
    EditorState.lineSeparator.of('\n'), // Prevent CM from treating carriage return as newline
    keymap.of({
        key: 'Shift-Enter',
        run: ({ state, dispatch }) => {
            dispatch(state.replaceSelection('\r'));
            return true;
        }
    })
];

let inputSequence = '';
export const insertChar = EditorView.domEventHandlers({
    'keydown': event => {
        if(event.altKey && event.key.match(/^[0-9a-f]$/i)) {
            inputSequence += event.key;
            event.preventDefault();
        }
    },
    'keyup': (event, view) => {
        if(event.key == 'Alt' && inputSequence) {
            try {
                const codepoint = parseInt(inputSequence, 16);
                if(codepoint != 0)
                    view.dispatch(view.state.replaceSelection(
                        String.fromCodePoint(codepoint)
                    ));
            }
            catch(e) {}
            inputSequence = '';
        }
    }
});

class UnprintableWidget extends WidgetType {
    constructor(value) {
        super();
        this.value = value;
    }
    toDOM() {
        return <span title={'\\u' + this.value.toString(16)}>•</span>
    }
}

const unprintableDecorator = new MatchDecorator({
    regexp: /[\x01-\x08\x0B-\x1F]/g,
    decoration: match => Decoration.replace({
        widget: new UnprintableWidget(match[0].charCodeAt(0))
    })
});

export const showUnprintables = ViewPlugin.fromClass(
    class {
        constructor(view) {
            this.decorations = unprintableDecorator.createDeco(view);
        }
        update(update) {
            this.decorations = unprintableDecorator.updateDeco(update, this.decorations);
        }
    },
    { decorations: plugin => plugin.decorations }
);
