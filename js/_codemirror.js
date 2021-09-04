import { EditorState }        from '@codemirror/state';
import { EditorView, keymap } from '@codemirror/view';

export { EditorState, EditorView };

// Extensions.
import { carriageReturn, insertChar, showUnprintables } from './_codemirror_unprintable.js';
import { insertTab, standardKeymap }                    from '@codemirror/commands';
import { lineNumbers }                                  from '@codemirror/gutter';
import { defaultHighlightStyle, HighlightStyle, tags }  from '@codemirror/highlight';
import { history, historyKeymap }                       from '@codemirror/history';
import { indentOnInput }                                from '@codemirror/language';
import { bracketMatching }                              from '@codemirror/matchbrackets';
import { StreamLanguage }                               from '@codemirror/stream-parser';
import { oneDarkTheme, oneDarkHighlightStyle }          from '@codemirror/theme-one-dark';

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
import { nim }         from 'nim-codemirror-mode';
import { pascal }      from '@codemirror/legacy-modes/mode/pascal';
import { perl }        from '@codemirror/legacy-modes/mode/perl';
import { php }         from '@codemirror/lang-php';
import { powerShell }  from '@codemirror/legacy-modes/mode/powershell';
import { python }      from '@codemirror/lang-python';
import { ruby }        from '@codemirror/legacy-modes/mode/ruby';
import { rust }        from '@codemirror/lang-rust';
import { shell }       from '@codemirror/legacy-modes/mode/shell';
import { sql, SQLite } from '@codemirror/lang-sql';
import { swift }       from '@codemirror/legacy-modes/mode/swift';

// For some reason, this doesn't fully work unless added to both themes.
const asmErrorTooltip = {
    '&:before':   { borderTopColor: 'var(--color)' },
    'background': 'var(--color)',
    'color':      'var(--background)',
};

const fontFamily =
    "'SFMono-Regular', Menlo, Consolas, 'Liberation Mono', Courier, monospace";

export const extensions = {
    // Extensions.
    base: [
        carriageReturn, defaultHighlightStyle, history(),
        indentOnInput(), insertChar, lineNumbers(), showUnprintables,
        keymap.of([
            { key: 'Tab', run: insertTab },
            ...historyKeymap, ...standardKeymap,
        ]),
        EditorView.theme({
            '.cm-asm-error': { textDecoration: 'underline var(--asm-error)' },
            '.cm-asm-error-tooltip':    asmErrorTooltip,
            '.cm-content':              { fontFamily },
            '.cm-gutters':              { fontFamily },
            '.cm-tooltip':              { fontFamily },
            '.cm-tooltip-autocomplete': { fontFamily },
        }, { dark: false }),
    ],
    bracketMatching: bracketMatching(),
    dark: [
        EditorView.theme({
            '&': { background: 'var(--background)', color: 'var(--color)' },
            ".cm-asm-dump":          { color: "var(--asm-dump)" },
            '.cm-asm-error-tooltip': asmErrorTooltip,
            '.cm-gutters':           { background: 'var(--light-grey)' },
        }, { dark: true }),
        HighlightStyle.define([
            { color: '#98c379', tag: tags.literal },
            { color: '#e06c75', tag: tags.regexp  },
        ]),
        oneDarkTheme,
        oneDarkHighlightStyle,
    ],
    defaultHighlightStyle,

    // Languages.
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
    'nim':        StreamLanguage.define(nim( {}, {} )),
    'pascal':     StreamLanguage.define(pascal),
    'perl':       StreamLanguage.define(perl),
    'php':        php({ plain: true }),
    'powershell': StreamLanguage.define(powerShell),
    'python':     python(),
    // TODO raku
    'ruby':       StreamLanguage.define(ruby),
    'rust':       rust(),
    'sql':        sql({ dialect: SQLite }),
    'swift':      StreamLanguage.define(swift),
    // TODO v
    // TODO zig
};
