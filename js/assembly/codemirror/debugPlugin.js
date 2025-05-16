import { EditorView, ViewPlugin, ViewUpdate, Decoration, hoverTooltip } from '@codemirror/view';
import { EditorState, ChangeSet } from '@codemirror/state';
import { ASMStateField } from './compilerPlugin.js';

var debugEnabled = false;

export const debugPlugin = [
    EditorView.baseTheme({
        '.red': { background: 'lightcoral' },
        '.blue': { background: 'lightblue' },
        '.cm-asm-debug-compiled .red, .red .cm-asm-debug-compiled': { background: 'indianred' },
        '.cm-asm-debug-compiled .blue, .blue .cm-asm-debug-compiled': { background: 'dodgerblue' },
        '.cm-asm-debug-compiled': { background: '#ddd' }
    }),
    hoverTooltip((view, pos) => {
        if(!debugEnabled)
            return null;
        const node = view.state.field(ASMStateField).head.find(pos);
        if(!node)
            return null;
        const instr = node.statement;
        return {
            pos: instr.range.start,
            end: Math.min(instr.range.end, view.state.doc.length),
            above: true,
            create: view => {
                let dom = document.createElement('div');
                dom.textContent = `${instr.constructor.name} (#${instr.id})`;
                dom.className = 'cm-asm-error-tooltip';
                return { dom };
            }
        }
    }),
    EditorView.domEventHandlers({
        mousedown: (event, view) => {
            if(debugEnabled && event.ctrlKey)
            {
                console.log(view.state.field(ASMStateField).head.find(view.posAtCoords(event)));
                return true;
            }
        },
        keydown: (event, view) => {
            if(event.key == 'F3')
            {
                debugEnabled = !debugEnabled;
                view.dispatch(ChangeSet.empty(0));
                return true;
            }
        }
    }),
    ViewPlugin.fromClass(
        class
        {
            /** @param {EditorView} view */
            constructor(view)
            {
                this.markInstructions(view.state);
            }

            /** @param {ViewUpdate} update */
            update(update)
            {
                this.markInstructions(update.state);
            }

            /** @param {EditorState} state */
            markInstructions(state)
            {
                if(!debugEnabled)
                {
                    this.decorations = Decoration.set([]);
                    return;
                }
                let instrMarks = [];
                let i = 0;
                state.field(ASMStateField).iterate(instr => {
                    instrMarks.push(Decoration.mark({
                        class: i++ % 2 ? 'blue' : 'red'
                    }).range(instr.range.start, instr.range.end))
                });
                let compiledRange = state.field(ASMStateField).compiledRange;
                if(compiledRange.length > 0)
                    instrMarks.push(Decoration.mark({
                        class: 'cm-asm-debug-compiled'
                    }).range(compiledRange.start, compiledRange.end));
                this.decorations = Decoration.set(instrMarks, true);
            }
        },
        { decorations: plugin => plugin.decorations }
    )
];