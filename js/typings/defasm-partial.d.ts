declare module '@defasm/core' {
    import { ASMError, AssemblyState as AssemblyStateInterface, Range, Statement, StatementNode } from "./defasm-interface.d.ts";

    export class AssemblyState implements AssemblyStateInterface {
        defaultSyntax: {
        intel?: boolean;
        prefix?: boolean;
        };
        bitness: number;
        head: StatementNode;
        compiledRange: Range;
        errors: ASMError[];

        compile(
        source: string,
        options?: {
            haltOnError?: boolean;
            range?: Range;
            doSecondPass?: boolean
        }
        ): void;

        secondPass(haltOnError?: boolean): void;

        iterate(callback: (instr: Statement, line: number) => void): void;

        bytesPerLine(callback: (bytes: Uint8Array, line: number) => void): void;

        constructor({
            syntax = {
                intel: false,
                prefix: true
            },
            writableText = false,
            bitness = 64
        } = {});
    }
}