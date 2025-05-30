export interface Message {
    type: 'error' | 'warning' | string;
    line?: number;
    endLine?: number;
    startColumn?: number;
    endColumn?: number;
    message: string;
}

export interface CodeFragment {
    line?: number;
    endLine?: number;
    startColumn?: number;
    endColumn?: number;
    section?: string;
    address: number;
    bytes: Uint8Array;
}

export interface ToolResult<T> {
    value?: T;
    messages: Message[];
}

export interface Result<T> {
    value?: T;
    error?: string;
}

export interface Assembler {
    assemble(source: string | Uint8Array): ToolResult<Uint8Array>;
    parseExecutable(executable: Uint8Array): Result<CodeFragment[]>;
}

export function assembleAndParse(assembler: Assembler, source: string): ToolResult<CodeFragment[]> {
    const {value, messages} = assembler.assemble(source);

    if (!value) {
        return { messages };
    }

    const parseResult = assembler.parseExecutable(value);
    if(parseResult.error) {
        messages.push({type: 'error', message: parseResult.error});
    }

    return {
        value: parseResult.value,
        messages
    };
}
