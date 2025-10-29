import * as intf from "../../typings/defasm-interface";
import { assembleAndParse, Assembler } from "../assembler-interface";
import { ResizableBuffer } from "../resizable-buffer";

class Range implements intf.Range
{
    _start: number;
    length: number

    constructor(start = 0, length = 0)
    {
        if(start < 0 || length < 0)
            throw `Invalid range ${start} to ${start + length}`;
        this._start = start;
        this.length = length;
    }

    /** @param {Number} pos */
    includes(pos: number)
    {
        return this.end >= pos && pos >= this.start;
    }

    /** @param {Range} end */
    until(end: intf.Range)
    {
        return new Range(this.start, end.end - this.start);
    }

    /** @param {string} text */
    slice(text: string)
    {
        return text.slice(this.start, this.end);
    }

    get start() { return this._start; }
    set start(val) { this._start = val; }

    get end() { return this.start + this.length; }
}

class ASMError implements intf.ASMError {
    constructor(public message: string, public range?: Range) {
    }
}

export class Section implements intf.Section {
    constructor(public name: string) {
    }
}

export class Statement implements intf.Statement {
  constructor(public range: Range, public section: Section, public bytes: Uint8Array, public length: number) {
  }
}

export class StatementNode implements intf.StatementNode
{
    statement: Statement | null;
    next: StatementNode | null;

    /** @param {Statement?} statement */
    constructor(statement: Statement | null = null)
    {
        this.statement = statement;
        /** @type {StatementNode?} */
        this.next = null;
    }

    find(pos: number): StatementNode | null
    {
        if(this.statement && this.statement.range.includes(pos))
            return this;
        return this.next ? this.next.find(pos) : null;
    }

    length()
    {
        let node: StatementNode | null = this;
        let length = 0;
        while(node)
        {
            if(node.statement)
                length += node.statement.length;
            node = node.next;
        }
        return length;
    }

    dump()
    {
        let node: StatementNode | null = this;
        let output, i = 0;

        output = new Uint8Array(this.length());

        while(node)
        {
            if(node.statement)
            {
                output.set(node.statement.bytes.subarray(0, node.statement.length), i);
                i += node.statement.length;
            }
            node = node.next;
        }

        return output;
    }
}

function determineLineRanges(text: string): Range[] {
    const ranges: Range[] = [];

    let currentStart = 0;

    for (let i = 0; i < text.length; i++) {
      const char = text[i];

      if (char === '\n') {
        const length = i - currentStart;
        ranges.push(new Range(currentStart, length));

        // Set the start of the next range
        currentStart = i + 1;
      }
    }

    // Add the last part
    ranges.push(new Range(currentStart, text.length - currentStart));

    return ranges;
}

export class AssemblyState implements intf.AssemblyState {
    constructor(private assembler: Assembler, public source = "", public errors: ASMError[] = [], public head = new StatementNode(), public compiledRange = new Range()) {
    }

    clone(): AssemblyState {
        return new AssemblyState(this.assembler, this.source, this.errors, this.head, this.compiledRange);
    }


    compile(source: string, {
        haltOnError = false,
        range: replacementRange = new Range(0, this.source.length),
        doSecondPass = true } = {}): void {
        this.source =
            /* If the given range is outside the current
            code's span, fill the in-between with newlines */
            this.source.slice(0, replacementRange.start).padEnd(replacementRange.start, '\n') +
            source +
            this.source.slice(replacementRange.end);

        // TODO: should we set it to the whole range since we are in fact reassembling the whole source code every time?
        this.compiledRange = new Range(replacementRange.start, source.length);

        if(doSecondPass)
            this.secondPass(haltOnError);
    }

    secondPass(haltOnError?: boolean): void {
        const lineRanges = determineLineRanges(this.source);
        const initialRange = new Range(0, 0);
        const finalRange = new Range(this.source.length, 0);
        const wholeRange = new Range(0, this.source.length);
        function getRange(range: Range, line?: number, endLine?: number, startColumn?: number, endColumn?: number) {
            if(line) {
                range = lineRanges[line - 1];
                if(endLine || startColumn || endColumn) {
                    const endLineRange = endLine ? lineRanges[endLine - 1] : range;
                    let start = range.start;
                    let end = endLineRange.end;
                    if(endColumn) {
                        end = endLineRange.start + endColumn - 1;
                    }
                    if(startColumn) {
                        start = range.start + startColumn - 1;
                    }
                    range = new Range(start, end - start);
                }
            }
            return range;
        }

        const {value, messages} = assembleAndParse(this.assembler, this.source);

        const head = new StatementNode();
        if(value) {
            let lastStatement = undefined;
            value.sort((a, b) => {
                if(a.line !== b.line) {
                    return (a.line || Number.POSITIVE_INFINITY) - (b.line || Number.POSITIVE_INFINITY);
                }
                if((a.startColumn || 1) !== (b.startColumn || 1)) {
                    return (a.startColumn || 1) - (b.startColumn || 1);
                }
                if(a.startColumn !== b.endColumn) {
                    return (a.endColumn || Number.POSITIVE_INFINITY) - (b.endColumn || Number.POSITIVE_INFINITY);
                }
                return 0;
            });
            for(const fragment of value) {
                const statement = new Statement(
                    getRange(finalRange, fragment.line, fragment.endLine, fragment.startColumn, fragment.endColumn),
                    new Section(fragment.section || ".text"),
                    fragment.bytes,
                    fragment.bytes.length,
                );
                const node = new StatementNode(statement);
                if(lastStatement) {
                    lastStatement.next = node;
                } else {
                    head.next = node;
                }
                lastStatement = node;
            }
        }
        this.head = head;

        const errors = [];
        for(const message of messages) {
            if(message.type === "error" || message.line) {
                const range = messages.length > 1 ? initialRange : wholeRange;
                errors.push(new ASMError((message.type != "error" ? (message.type + ": ") : "") + message.message, getRange(range,message.line, message.endLine, message.startColumn, message.endColumn)));
            }
        }
        this.errors = errors;
    }

    iterate(func: (instr: intf.Statement, line: number) => void): void
    {
        let line = 1, nextLine = 0, node = this.head.next;
        while(nextLine != Infinity)
        {
            nextLine = this.source.indexOf('\n', nextLine) + 1 || Infinity;
            while(node && node.statement!.range.end < nextLine)
            {
                func(node.statement!, line);
                node = node.next;
            }
            line++;
        }
    }

    bytesPerLine(func: (bytes: Uint8Array, line: number) => void): void {
        let line = 1, nextLine = 0, node = this.head.next;
        while(nextLine != Infinity)
        {
            let bytes = new ResizableBuffer();
            nextLine = this.source.indexOf('\n', nextLine) + 1 || Infinity;
            while(node && node.statement!.range.start < nextLine)
            {
                let instr = node.statement!;
                bytes.append(instr.bytes, 0, instr.length);
                node = node.next;
            }
            func(bytes.array, line);
            line++;
        }
    }
}
