import { EditorState }        from '@codemirror/state';
import { EditorView, keymap } from '@codemirror/view';

// Extensions.
import { closeBrackets, closeBracketsKeymap }          from '@codemirror/closebrackets';
import { defaultTabBinding, standardKeymap }           from '@codemirror/commands';
import { lineNumbers }                                 from '@codemirror/gutter';
import { defaultHighlightStyle, HighlightStyle, tags } from '@codemirror/highlight';
import { history, historyKeymap }                      from '@codemirror/history';
import { indentOnInput }                               from '@codemirror/language';
import { bracketMatching }                             from '@codemirror/matchbrackets';
import { StreamLanguage }                              from '@codemirror/stream-parser';
import { oneDarkTheme, oneDarkHighlightStyle }         from '@codemirror/theme-one-dark';

// Languages.
import { assembly }    from '@defasm/codemirror';
import { brainfuck }   from '@codemirror/legacy-modes/mode/brainfuck';
import { c, csharp }   from '@codemirror/legacy-modes/mode/clike';
import { cobol }       from '@codemirror/legacy-modes/mode/cobol';
import { commonLisp }  from '@codemirror/legacy-modes/mode/commonlisp';
import { crystal }     from '@codemirror/legacy-modes/mode/crystal';
import { diff }        from '@codemirror/legacy-modes/mode/diff';
import { fortran }     from '@codemirror/legacy-modes/mode/fortran';
import { fSharp }      from '@codemirror/legacy-modes/mode/mllike';
import { go }          from '@codemirror/legacy-modes/mode/go';
import { haskell }     from '@codemirror/legacy-modes/mode/haskell';
import { java }        from '@codemirror/lang-java';
import { javascript }  from '@codemirror/lang-javascript';
import { julia }       from '@codemirror/legacy-modes/mode/julia';
import { lua }         from '@codemirror/legacy-modes/mode/lua';
import { perl }        from '@codemirror/legacy-modes/mode/perl';
import { powerShell }  from '@codemirror/legacy-modes/mode/powershell';
import { python }      from '@codemirror/lang-python';
import { ruby }        from '@codemirror/legacy-modes/mode/ruby';
import { rust }        from '@codemirror/lang-rust';
import { shell }       from '@codemirror/legacy-modes/mode/shell';
import { sql, SQLite } from '@codemirror/lang-sql';
import { swift }       from '@codemirror/legacy-modes/mode/swift';

// Other.
import LZString from 'lz-string';

export { EditorState, EditorView, LZString };

// For some reason, this doesn't fully work unless added to both themes.
const asmErrorTooltip = {
    color: "var(--background)",
    backgroundColor: "var(--color)",
    "&:before": {
        borderTopColor: "var(--color)",
    }
};

export const darkThemeExtensions = [
    EditorView.theme({
        "&": {
            backgroundColor: "var(--background)",
            color: "var(--color)",
        },
        ".cm-gutters": {
            backgroundColor: "var(--light-grey)",
        },
        ".cm-asm-error-tooltip": asmErrorTooltip,
    }, { dark: true }),
    HighlightStyle.define([{
        tag: tags.literal,
        color: "#98c379",
    },
    {
        tag: tags.regexp,
        color: "#e06c75"
    }]),
    oneDarkTheme,
    oneDarkHighlightStyle,
];

const fontFamily = "'SFMono-Regular', Menlo, Consolas, 'Liberation Mono', Courier, monospace";

const defaultTheme = EditorView.theme({
    '.cm-content':              { fontFamily },
    '.cm-gutters':              { fontFamily },
    '.cm-tooltip':              { fontFamily },
    '.cm-tooltip-autocomplete': { fontFamily },
    '.cm-asm-error':            { textDecoration: 'underline var(--asm-error)' },
    '.cm-asm-error-tooltip':    asmErrorTooltip,
}, { dark: false });

export const defaultExtensions = [
    EditorView.lineWrapping,
    bracketMatching(),
    closeBrackets(),
    defaultHighlightStyle,
    defaultTheme,
    history(),
    indentOnInput(),
    keymap.of([
        ...closeBracketsKeymap, defaultTabBinding,
        ...historyKeymap, ...standardKeymap,
    ]),
    lineNumbers(),
];

export const languages = {
    'assembly':   assembly(),
    'bash':       StreamLanguage.define(shell),
    'brainfuck':  StreamLanguage.define(brainfuck),
    'c':          StreamLanguage.define(c),
    'c-sharp':    StreamLanguage.define(csharp),
    'cobol':      StreamLanguage.define(cobol),
    'crystal':    StreamLanguage.define(crystal),
    'diff':       StreamLanguage.define(diff),
    'f-sharp':    StreamLanguage.define(fSharp),
    // TODO fish
    'fortran':    StreamLanguage.define(fortran),
    'go':         StreamLanguage.define(go),
    'haskell':    StreamLanguage.define(haskell),
    // TODO hexagony
    // TODO j
    'java':       java(),
    'javascript': javascript(),
    'julia':      StreamLanguage.define(julia),
    'lisp':       StreamLanguage.define(commonLisp),
    'lua':        StreamLanguage.define(lua),
    // TODO nim
    'perl':       StreamLanguage.define(perl),
    // TODO php
    'powershell': StreamLanguage.define(powerShell),
    'python':     python(),
    // TODO raku
    'ruby':       StreamLanguage.define(ruby),
    'rust':       rust(),
    'sql':        sql({ dialect: SQLite }),
    'swift':      StreamLanguage.define(swift),
    // TODO v
    // TODO zig
}
