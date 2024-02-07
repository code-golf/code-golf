import { EditorState }                                    from '@codemirror/state';
import { EditorView, keymap, lineNumbers, drawSelection,
    highlightWhitespace } from '@codemirror/view';
import { $ }                                              from './_util';

export { EditorState, EditorView };

// Extensions.
import { carriageReturn, insertChar,
    showUnprintables }                           from './_codemirror_unprintable';
import { history, historyKeymap, indentLess, insertNewline,
    insertTab, standardKeymap, toggleComment }   from '@codemirror/commands';
import { tags }                                  from '@lezer/highlight';
import { bracketMatching, defaultHighlightStyle,
    HighlightStyle, LanguageSupport, StreamLanguage,
    syntaxHighlighting }                         from '@codemirror/language';
import { oneDarkTheme, oneDarkHighlightStyle }   from '@codemirror/theme-one-dark';
import { vim }                                   from '@replit/codemirror-vim';

// Languages.
import { assembly }           from '@defasm/codemirror';
import { brainfuck }          from 'codemirror-lang-brainfuck';
import { c, csharp, dart }    from './vendor/codemirror-clike';
import { clojure }            from '@codemirror/legacy-modes/mode/clojure';
import { cobol }              from './vendor/codemirror-cobol';
import { commonLisp }         from '@codemirror/legacy-modes/mode/commonlisp';
import { cpp }                from '@codemirror/lang-cpp';
import { crystal }            from '@codemirror/legacy-modes/mode/crystal';
import { d }                  from '@codemirror/legacy-modes/mode/d';
import { elixirLanguage }     from 'codemirror-lang-elixir';
import { factor }             from '@codemirror/legacy-modes/mode/factor';
import { forth }              from '@codemirror/legacy-modes/mode/forth';
import { fortran }            from '@codemirror/legacy-modes/mode/fortran';
import { fSharp, oCaml }      from '@codemirror/legacy-modes/mode/mllike';
import { go }                 from '@codemirror/legacy-modes/mode/go';
import { golfScript }         from 'codemirror-lang-golfscript';
import { haskell }            from '@codemirror/legacy-modes/mode/haskell';
import { j }                  from 'codemirror-lang-j';
import { janet }              from 'codemirror-lang-janet';
import { java }               from '@codemirror/lang-java';
import { javascriptLanguage } from '@codemirror/lang-javascript';
import { julia }              from '@codemirror/legacy-modes/mode/julia';
import { k }                  from 'codemirror-lang-k';
import { lua }                from '@codemirror/legacy-modes/mode/lua';
import { nim }                from 'nim-codemirror-mode';
import { pascal }             from '@codemirror/legacy-modes/mode/pascal';
import { perl }               from '@codemirror/legacy-modes/mode/perl';
import { phpLanguage }        from '@codemirror/lang-php';
import { powerShell }         from '@codemirror/legacy-modes/mode/powershell';
import { prolog }             from 'codemirror-lang-prolog';
import { pythonLanguage }     from '@codemirror/lang-python';
import { r }                  from '@codemirror/legacy-modes/mode/r';
import { raku }               from './vendor/codemirror-raku';
import { ruby }               from '@codemirror/legacy-modes/mode/ruby';
import { rust }               from '@codemirror/lang-rust';
import { shell }              from '@codemirror/legacy-modes/mode/shell';
import { SQLite }             from '@codemirror/lang-sql';
import { swift }              from '@codemirror/legacy-modes/mode/swift';
import { tcl }                from '@codemirror/legacy-modes/mode/tcl';
import { stex }               from '@codemirror/legacy-modes/mode/stex';
import { wrenLanguage }       from '@exercism/codemirror-lang-wren';
import { zig }                from 'codemirror-lang-zig';

// Bypass default constructors so we only get highlighters and not extensions.
const elixir     = new LanguageSupport(elixirLanguage);
const javascript = new LanguageSupport(javascriptLanguage);
const php        = new LanguageSupport(phpLanguage.configure({ top: 'Program' }));
const python     = new LanguageSupport(pythonLanguage);
const sql        = new LanguageSupport(SQLite.language);
const wren       = new LanguageSupport(wrenLanguage);

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
        carriageReturn, showUnprintables,        
        drawSelection(),
        highlightWhitespace(),
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
    'editor': [
        history(),
        insertChar,
        keymap.of([
            // Replace "enter" with a non auto indenting action.
            ...historyKeymap, ...standardKeymap.filter(k => k.key != 'Enter'),
            { key: 'Enter', run: insertNewline },
            { key: 'Tab',   run: insertTab, shift: indentLess },
            { key: 'Mod-/', run: toggleComment },
        ]),
        lineNumbers()
    ],
    'bracketMatching': bracketMatching(),
    'vim': vim({ status: true }),

    // Languages.
    'assembly':   assembly(),
    'bash':       StreamLanguage.define(shell),
    // TODO basic
    'brainfuck':  brainfuck(),
    'c':          StreamLanguage.define(c),
    'c-sharp':    StreamLanguage.define(csharp),
    'clojure':    StreamLanguage.define(clojure),
    'cobol':      StreamLanguage.define(cobol),
    'coconut':    python,
    'cpp':        cpp(),
    'crystal':    StreamLanguage.define(crystal),
    'd':          StreamLanguage.define(d),
    'dart':       StreamLanguage.define(dart),
    'elixir':     elixir,
    'f-sharp':    StreamLanguage.define(fSharp),
    'factor':     StreamLanguage.define(factor),
    // TODO fish
    'forth':      StreamLanguage.define({ ...forth, languageData: { commentTokens: { line: '\\' } } }),
    'fortran':    StreamLanguage.define({ ...fortran, languageData: { commentTokens: { line: '!' } } }),
    'go':         StreamLanguage.define(go),
    'golfscript': golfScript(),
    'haskell':    StreamLanguage.define(haskell),
    // TODO hexagony
    'j':          j(),
    'janet':      janet(),
    'java':       java(),
    'javascript': javascript,
    'julia':      StreamLanguage.define(julia),
    'k':          k(),
    'lisp':       StreamLanguage.define(commonLisp),
    'lua':        StreamLanguage.define(lua),
    'nim':        StreamLanguage.define({ ...nim( {}, {} ), languageData: { commentTokens: { line: '#' } } }),
    'ocaml':      StreamLanguage.define(oCaml),
    'pascal':     StreamLanguage.define(pascal),
    'perl':       StreamLanguage.define(perl),
    'php':        php,
    'powershell': StreamLanguage.define(powerShell),
    'prolog':     prolog(),
    'python':     python,
    'r':          StreamLanguage.define(r),
    'raku':       StreamLanguage.define(raku),
    'ruby':       StreamLanguage.define(ruby),
    'rust':       rust(),
    // TODO sed
    'sql':        sql,
    'swift':      StreamLanguage.define(swift),
    'tcl':        StreamLanguage.define(tcl),
    'tex':        StreamLanguage.define(stex),
    // TODO v
    // TODO viml
    'wren':       wren,
    'zig':        zig(),
};

// Order matters, unshift the dark stuff onto the front.
if (matchMedia(JSON.parse($('#dark-mode-media-query').innerText)).matches)
    extensions.base.unshift(
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
    );
