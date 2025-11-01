// Type definitions for the subset of DefAsm used by the CodeMirror extension
export interface Range {
  length: number;
  start: number;
  readonly end: number;
  includes(pos: number): boolean;
  slice(text: string): string;
  until(end: Range): Range;
}

export interface Section {
  name: string;
}

export interface Statement {
  range: Range;
  section: Section;
  length: number;
  bytes: Uint8Array;
  syntax?: { intel?: boolean; prefix?: boolean };
}

export interface StatementNode {
  statement: Statement | null;
  next: StatementNode | null;
  find(pos: number): StatementNode | null;
  length(): number;
  dump(): Uint8Array;
}

export interface ASMError {
  message: string;
  range?: Range;
}

export interface AssemblyState {
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
}
