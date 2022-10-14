import { EditorState }                     from '@codemirror/state';
import { EditorView, keymap, lineNumbers,
    drawSelection }                        from '@codemirror/view';

export { EditorState, EditorView };

// Extensions.
import { carriageReturn, insertChar,
    showUnprintables }                           from './_codemirror_unprintable';
import { history, historyKeymap, insertNewline,
    insertTab, standardKeymap, toggleComment }   from '@codemirror/commands';
import { tags }                                  from '@lezer/highlight';
import { bracketMatching, defaultHighlightStyle,
    HighlightStyle, StreamLanguage,
    syntaxHighlighting }                         from '@codemirror/language';
import { oneDarkTheme, oneDarkHighlightStyle }   from '@codemirror/theme-one-dark';
import { vim }                                   from '@replit/codemirror-vim';

// Languages.
import { assembly }        from '@defasm/codemirror';
import { brainfuck }       from 'codemirror-lang-brainfuck';
import { c, csharp, dart } from '@codemirror/legacy-modes/mode/clike';
import { cobol }           from './vendor/codemirror-cobol';
import { commonLisp }      from '@codemirror/legacy-modes/mode/commonlisp';
import { cpp }             from '@codemirror/lang-cpp';
import { crystal }         from '@codemirror/legacy-modes/mode/crystal';
import { d }               from '@codemirror/legacy-modes/mode/d';
import { fortran }         from '@codemirror/legacy-modes/mode/fortran';
import { fSharp }          from '@codemirror/legacy-modes/mode/mllike';
import { go }              from '@codemirror/legacy-modes/mode/go';
import { golfScript }      from 'codemirror-lang-golfscript';
import { haskell }         from '@codemirror/legacy-modes/mode/haskell';
import { j }               from 'codemirror-lang-j';
import { java }            from '@codemirror/lang-java';
import { javascript }      from '@codemirror/lang-javascript';
import { julia }           from '@codemirror/legacy-modes/mode/julia';
import { k }               from 'codemirror-lang-k';
import { lua }             from '@codemirror/legacy-modes/mode/lua';
import { nim }             from 'nim-codemirror-mode';
import { pascal }          from '@codemirror/legacy-modes/mode/pascal';
import { perl }            from '@codemirror/legacy-modes/mode/perl';
import { php }             from '@codemirror/lang-php';
import { powerShell }      from '@codemirror/legacy-modes/mode/powershell';
import { python }          from '@codemirror/lang-python';
import { raku }            from './vendor/codemirror-raku';
import { ruby }            from '@codemirror/legacy-modes/mode/ruby';
import { rust }            from '@codemirror/lang-rust';
import { shell }           from '@codemirror/legacy-modes/mode/shell';
import { sql, SQLite }     from '@codemirror/lang-sql';
import { swift }           from '@codemirror/legacy-modes/mode/swift';
import { tcl }             from '@codemirror/legacy-modes/mode/tcl';
import { wren }            from '@exercism/codemirror-lang-wren';

// For some reason, this doesn't fully work unless added to both themes.
const asmErrorTooltip = {
    '&:before':   { borderTopColor: 'var(--color)' },
    'background': 'var(--color)',
    'color':      'var(--background)',
};

const fontFamily = "'Source Code Pro', monospace";

export const extensions = {
    // Extensions.
    'base': [
        carriageReturn, history(), insertChar, lineNumbers(), showUnprintables,
        keymap.of([
            // Replace "enter" with a non auto indenting action.
            ...historyKeymap, ...standardKeymap.filter(k => k.key != 'Enter'),
            { key: 'Enter', run: insertNewline },
            { key: 'Tab',   run: insertTab },
            { key: 'Mod-/', run: toggleComment },
        ]),
        drawSelection(),
        syntaxHighlighting(defaultHighlightStyle),
        EditorView.theme({
            '.cm-asm-error': { textDecoration: 'underline var(--asm-error)' },
            '.cm-asm-error-tooltip':    asmErrorTooltip,
            '.cm-content':              { fontFamily },
            '.cm-gutters':              { fontFamily },
            '.cm-tooltip':              { fontFamily },
            '.cm-tooltip-autocomplete': { fontFamily },
        }, { dark: false }),
    ],
    'bracketMatching': bracketMatching(),
    'dark': [
        EditorView.theme({
            '&': { background: 'var(--background)', color: 'var(--color)' },
            '.cm-asm-dump':               { color: 'var(--asm-dump)' },
            '.cm-asm-error-tooltip':      asmErrorTooltip,
            '.cm-gutters':                { background: 'var(--light-grey)' },
            '.cm-content':                { caretColor: 'var(--color)' },
            '.cm-cursor, .cm-dropCursor': { borderLeftColor: 'var(--color)' },
        }, { dark: true }),
        syntaxHighlighting(HighlightStyle.define([
            { color: '#98c379', tag: tags.literal },
            { color: '#e06c75', tag: tags.regexp  },
        ])),
        oneDarkTheme,
        syntaxHighlighting(oneDarkHighlightStyle),
    ],
    'vim': vim({ status: true }),

    // Languages.
    'assembly':   assembly(),
    'bash':       StreamLanguage.define(shell),
    'brainfuck':  brainfuck(),
    'c':          StreamLanguage.define(c),
    'c-sharp':    StreamLanguage.define(csharp),
    'cobol':      StreamLanguage.define(cobol),
    'cpp':        cpp(),
    'crystal':    StreamLanguage.define(crystal),
    'd':          StreamLanguage.define(d),
    'dart':       StreamLanguage.define(dart),
    // TODO elixir
    'f-sharp':    StreamLanguage.define(fSharp),
    // TODO fish
    'fortran':    StreamLanguage.define({ ...fortran, languageData: { commentTokens: { line: '!' } } }),
    'go':         StreamLanguage.define(go),
    'golfscript': golfScript(),
    'haskell':    StreamLanguage.define(haskell),
    // TODO hexagony
    'j':          j(),
    'java':       java(),
    'javascript': javascript(),
    'julia':      StreamLanguage.define(julia),
    'k':          k(),
    'lisp':       StreamLanguage.define(commonLisp),
    'lua':        StreamLanguage.define(lua),
    'nim':        StreamLanguage.define({ ...nim( {}, {} ), languageData: { commentTokens: { line: '#' } } }),
    'pascal':     StreamLanguage.define(pascal),
    'perl':       StreamLanguage.define(perl),
    'php':        php({ plain: true }),
    'powershell': StreamLanguage.define(powerShell),
    // TODO prolog
    'python':     python(),
    'raku':       StreamLanguage.define(raku),
    'ruby':       StreamLanguage.define(ruby),
    'rust':       rust(),
    // TODO sed
    'sql':        sql({ dialect: SQLite }),
    'swift':      StreamLanguage.define(swift),
    'tcl':        StreamLanguage.define(tcl),
    // TODO v
    // TODO viml
    'wren':       wren(),
    // TODO zig
};
