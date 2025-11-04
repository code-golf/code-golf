import { AssemblyState } from "../../typings/defasm-partial";

/**
 * Create a CodeMirror extension for assembly language
 * @param config Configuration options for the assembly extension
 */
export function assemblyLanguage(
  config?: {
    bitness?: number;
    defaultSyntax?: {
      intel?: boolean;
      prefix?: boolean;
    };
  }
): Extension;

/**
 * Create a CodeMirror extension for assembly using a DefAsm-compatible assembler
 * @param assemblyState The assembly state. Must be compatible with DefAsm's AssemblyState
 * @param config Configuration options for the assembly extension
 */
export function assemblyIde(
  assemblyState: AssemblyState,
  config?: {
    byteDumps?: boolean;
    debug?: boolean;
    errorMarking?: boolean;
    errorTooltips?: boolean;
  }
): Extension;

export const ASMStateField: StateField<AssemblyState>;