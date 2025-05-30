import { EditorView, ViewPlugin, Decoration, WidgetType, hoverTooltip } from '@codemirror/view';
import { ASMStateField } from "./compilerPlugin";
import { EditorState } from '@codemirror/state';

class EOLError extends WidgetType
{
    constructor()
    {
        super();
    }

    toDOM()
    {
        let node = document.createElement('span');
        node.setAttribute('aria-hidden', 'true');
        node.className = 'cm-asm-error';
        node.style.position = 'absolute';
        node.innerText = ' ';
        return node;
    }
}

export const errorMarker = [
    EditorView.baseTheme({
        '.cm-asm-error': {
            textDecoration: "underline red"
        }
    }),
    ViewPlugin.fromClass(
        class
        {
            constructor(view) { this.markErrors(view.state); }
            update(update) { if(update.docChanged) this.markErrors(update.state); }

            /** @param {EditorState} state */
            markErrors(state)
            {
                this.marks = Decoration.set(state.field(ASMStateField).errors.map(error => {
                    let content = state.sliceDoc(error.range.start, error.range.end);
                    if(content == '\n' || !content)
                        return Decoration.widget({
                            widget: new EOLError(),
                            side: 1
                        }).range(error.range.start);
                        
                    return Decoration.mark({
                        class: 'cm-asm-error'
                    }).range(error.range.start, error.range.end);
                }));
            }
        },
        { decorations: plugin => plugin.marks }
    )
];

export const errorTooltipper = [
    EditorView.baseTheme({
        '.cm-asm-error-tooltip': {
            fontFamily: "monospace",
            borderRadius: ".25em",
            padding: ".1em .25em",
            color: "#eee",
            backgroundColor: "black !important",
            "&:before": {
                position: "absolute",
                content: '""',
                left: ".3em",
                marginLeft: "-.1em",
                bottom: "-.3em",
                borderLeft: ".3em solid transparent",
                borderRight: ".3em solid transparent",
                borderTop: ".3em solid black"
            }
        },
        '&dark .cm-asm-error-tooltip': {
            color: "black",
            backgroundColor: "#eee !important",
            "&:before": {
                borderTop: ".3em solid #eee"
            }
        }
    }),
    hoverTooltip((view, pos) => {
        for(let { range, message } of view.state.field(ASMStateField).errors)
            if(range.start <= pos && range.end >= pos)
                return {
                    pos: range.start,
                    end: Math.min(range.end, view.state.doc.length),
                    above: true,
                    create: view => {
                        let dom = document.createElement('div');
                        dom.textContent = message;
                        dom.className = 'cm-asm-error-tooltip';
                        return { dom };
                    }
                }

        return null;
    })
];