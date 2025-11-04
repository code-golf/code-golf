import { EditorView, ViewPlugin, ViewUpdate, Decoration, WidgetType } from '@codemirror/view';
import { EditorState, StateField, Facet, RangeValue, RangeSet, RangeSetBuilder, StateEffect } from '@codemirror/state';

class Range
{
    constructor(start = 0, length = 0)
    {
        if(start < 0 || length < 0)
            throw `Invalid range ${start} to ${start + length}`;
        this._start = start;
        this.length = length;
    }

    /** @param {Number} pos */
    includes(pos)
    {
        return this.end >= pos && pos >= this.start;
    }

    /** @param {Range} end */
    until(end)
    {
        return new Range(this.start, end.end - this.start);
    }

    /** @param {string} text */
    slice(text)
    {
        return text.slice(this.start, this.end);
    }

    get start() { return this._start; }
    set start(val) { this._start = val; }

    get end() { return this.start + this.length; }
}

/** @type {StateField<AssemblyState>} */
export const ASMStateField = StateField.define({
    create: state => {
        throw new Exception("You must use ASMStateField.init to create ASMStateField fields");
    },
    update: (state, transaction) => {
        if(!transaction.docChanged)
            return state;

        // XXX: TODO: CodeMirror requires the state field to be immutable, so this is BROKEN with DefAsm and any compatible implementation that doesn't support clone()
        // we do this because it's the way the code used to work, and it looks like it does somehow usually work in practice
        if(state.clone)
            state = state.clone();
        /* In case there are multiple changed ranges, we'll compile each
        range separately and only run the second pass on the final state. */
        let offset = 0;
        transaction.changes.iterChanges(
            (fromA, toA, fromB, toB) => {
                state.compile(transaction.state.sliceDoc(fromB, toB), { range: new Range(fromA + offset, toA - fromA), doSecondPass: false });
                offset += toB - fromB - (toA - fromA);
            }
        );

        state.secondPass();
        return state;
    }
});

export const ASMLanguageData = EditorState.languageData.of((state, pos, side) => {
    pos = state.doc.lineAt(pos).from;

    var asmState = state.field(ASMStateField);
    var lastInstr = null;

    asmState.iterate(instr => {
        if(instr.range.start < pos)
            lastInstr = instr;
    });
    return [{
        commentTokens: { line:
            (lastInstr ? lastInstr.syntax : asmState.defaultSyntax)?.intel ? ';' : '#' }
    }];
})

export class ASMColor extends RangeValue
{
    /** @param {String} color */
    constructor(color)
    {
        super();
        this.color = color;
    }
    eq(other)
    {
        return this.color == other.color;
    }
}

/** @type {Facet<RangeSet<ASMColor>>} */
export const ASMColorFacet = Facet.define();

export const SectionColors = ASMColorFacet.compute(
    ['doc'],
    state => {
        let assemblyState = state.field(ASMStateField), offset = 0;
        /** @type {RangeSetBuilder<ASMColor>} */
        let builder = new RangeSetBuilder();
        assemblyState.iterate((instr, line) => {
            let sectionName = instr.section.name;
            let color = sectionName == ".text"   ? "#666" :
                        sectionName == ".data"   ? "#66A" :
                        sectionName == ".bss"    ? "#6A6" :
                        sectionName == ".rodata" ? "#AA6" :
                                                   "#A66";
            builder.add(offset, offset += instr.length, new ASMColor(color));
        });
        return builder.finish();
    }
);

class AsmDumpWidget extends WidgetType
{
    /**
     * @param {Uint8Array} bytes
     * @param {Number} offset
     * @param {RangeSet<ASMColor>} colors
     */
    constructor(bytes, offset, colors)
    {
        super();
        this.bytes = bytes;
        this.offset = offset;
        this.colors = colors;
    }

    toDOM()
    {
        let node = document.createElement('span');
        node.setAttribute('aria-hidden', 'true');
        node.className = 'cm-asm-dump';
        node.style.marginLeft = this.offset + 'px';
        let colorCursor = this.colors.iter();

        for(let i = 0; i < this.bytes.length; i++)
        {
            let text = this.bytes[i].toString(16).toUpperCase().padStart(2, '0');
            let span = document.createElement('span');

            while(colorCursor.to <= i)
                colorCursor.next();

            if(colorCursor.from <= i && i < colorCursor.to)
                span.style.color = colorCursor.value.color;

            span.innerText = text;
            node.appendChild(span);
        }

        return node;
    }

    eq(widget)
    {
        if(this.offset != widget.offset || this.bytes.length != widget.bytes.length)
            return false;

        for(let i = 0; i < this.bytes.length; i++)
            if(this.bytes[i] != widget.bytes[i])
                return false;

        // RangeSet.eq doesn't work for some reason
        let oldCursor = widget.colors.iter(), newCursor = this.colors.iter();
        while(true)
        {
            if(newCursor.value === null)
            {
                if(oldCursor.value === null)
                    break;
                return false;
            }

            if(!newCursor.value.eq(oldCursor.value) ||
                newCursor.from !== oldCursor.from ||
                newCursor.to !== oldCursor.to)
                return false;

            oldCursor.next(); newCursor.next();
        }

        return true;
    }

    /** @param {HTMLElement} node */
    updateDOM(node)
    {
        node.style.marginLeft = this.offset + 'px';
        let colorCursor = this.colors.iter();

        for(let i = 0; i < this.bytes.length; i++)
        {
            while(colorCursor.to <= i)
                colorCursor.next();

            let text = this.bytes[i].toString(16).toUpperCase().padStart(2, '0');
            if(i < node.childElementCount)
            {
                let span = node.children.item(i);
                if(span.innerText !== text)
                    span.innerText = text;
                if(colorCursor.from <= i && i < colorCursor.to)
                {
                    if(colorCursor.value.color !== span.style.color)
                        span.style.color = colorCursor.value.color;
                }
                else
                    span.style.color = "";
            }
            else
            {
                let span = document.createElement('span');
                if(colorCursor.value !== null)
                    span.style.color = colorCursor.value.color;

                span.innerText = text;
                node.appendChild(span);
            }
            while(colorCursor.to < i)
                colorCursor.next();
        }
        while(node.childElementCount > this.bytes.length)
            node.removeChild(node.lastChild);

        return true;
    }
}

/* Convert tabs to spaces, for proper width measurement */
function expandTabs(text, tabSize)
{
    let result = "", i = tabSize;
    for(let char of text)
    {
        if(char == '\t')
        {
            result += ' '.repeat(i);
            i = tabSize;
        }
        else
        {
            result += char;
            i = i - 1 || tabSize;
        }
    }
    return result;
}

export const ASMFlush = StateEffect.define();

export const byteDumper = [
    EditorView.baseTheme({
        '.cm-asm-dump'       : { fontStyle: "italic" },
        '.cm-asm-dump > span': { marginRight: "1ch" },
        '&dark .cm-asm-dump' : { color: "#AAA" }
    }),
    ViewPlugin.fromClass(class {
        /** @param {EditorView} view */
        constructor(view)
        {
            this.ctx        = document.createElement('canvas').getContext('2d');
            this.lineWidths = [];

            this.decorations = Decoration.set([]);

            // This timeout is required to let the content DOM's style be calculated
            setTimeout(() => {
                let style = window.getComputedStyle(view.contentDOM);

                this.ctx.font = `${
                    style.getPropertyValue('font-style')
                } ${
                    style.getPropertyValue('font-variant')
                } ${
                    style.getPropertyValue('font-weight')
                } ${
                    style.getPropertyValue('font-size')
                } ${
                    style.getPropertyValue('font-family')
                }`;
                
                
                this.updateWidths(0, view.state.doc.length, 0, view.state);
                this.makeAsmDecorations(view.state);
                view.dispatch();
            }, 1);
        }

        /** @param {ViewUpdate} update */
        update(update)
        {
            if(!update.docChanged && !update.transactions.some(
                tr => tr.effects.some(
                    effect => effect.is(ASMFlush)
                )
            ))
                return;

            let state = update.view.state;

            update.changes.iterChangedRanges(
                (fromA, toA, fromB, toB) => {
                    let removedLines =
                        update.startState.doc.lineAt(toA).number
                        -
                        update.startState.doc.lineAt(fromA).number;
                    this.updateWidths(fromB, toB, removedLines, state);
                }
            );

            this.makeAsmDecorations(update.state);
        }

        updateWidths(from, to, removedLines, { doc, tabSize })
        {
            let start = doc.lineAt(from).number;
            let end   = doc.lineAt(to).number;
            let newWidths = [];
            
            for(let i = start; i <= end; i++)
                newWidths.push(this.ctx.measureText(expandTabs(doc.line(i).text, tabSize)).width);
            
            this.lineWidths.splice(start - 1, removedLines + 1, ...newWidths);
        }

        /** @param {EditorState} state */
        makeAsmDecorations(state)
        {
            let doc       = state.doc;
            let maxOffset = Math.max(...this.lineWidths) + 50;
            let widgets   = [];

            let asmColors = state.facet(ASMColorFacet);
            let assemblyState = state.field(ASMStateField);
            let i = 0;

            assemblyState.bytesPerLine((bytes, line) => {
                if(bytes.length > 0)
                {
                    /** @type {RangeSetBuilder<ASMColor>} */
                    let builder = new RangeSetBuilder();
                    RangeSet.spans(asmColors, i, i + bytes.length, {
                        span: (from, to, active) => {
                            if(active.length > 0)
                                builder.add(from - i, to - i, active[active.length - 1]);
                        }
                    });
                    let colors = builder.finish();
                    i += bytes.length;

                    widgets.push(Decoration.widget({
                            widget: new AsmDumpWidget(
                                bytes,
                                maxOffset - this.lineWidths[line - 1],
                                colors
                            ), side: 2
                        }).range(doc.line(line).to));
                }
            });

            this.decorations = Decoration.set(widgets);
        }

    }, { decorations: plugin => plugin.decorations })
];
