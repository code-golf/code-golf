import { EditorState }                                    from '@codemirror/state';
import { EditorView, keymap, lineNumbers, drawSelection,
    highlightWhitespace } from '@codemirror/view';
import { $ }                                              from './_util';

export { EditorState, EditorView };

// Extensions.
import { carriageReturn, insertChar, insertCharState,
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
import { apl }                                      from '@codemirror/legacy-modes/mode/apl';
import { assembly }                                 from '@defasm/codemirror';
import { brainfuck }                                from 'codemirror-lang-brainfuck';
import { c, csharp, dart, kotlin, scala, squirrel } from './vendor/codemirror-clike';
import { clojure }                                  from '@codemirror/legacy-modes/mode/clojure';
import { cobol }                                    from './vendor/codemirror-cobol';
import { coffeeScript }                             from '@codemirror/legacy-modes/mode/coffeescript';
import { commonLisp }                               from '@codemirror/legacy-modes/mode/commonlisp';
import { cpp }                                      from '@codemirror/lang-cpp';
import { crystal }                                  from '@codemirror/legacy-modes/mode/crystal';
import { d }                                        from '@codemirror/legacy-modes/mode/d';
import { elixirLanguage }                           from 'codemirror-lang-elixir';
import { erlang }                                   from '@codemirror/legacy-modes/mode/erlang';
import { factor }                                   from '@codemirror/legacy-modes/mode/factor';
import { forth }                                    from '@codemirror/legacy-modes/mode/forth';
import { fortran }                                  from '@codemirror/legacy-modes/mode/fortran';
import { fSharp, oCaml }                            from '@codemirror/legacy-modes/mode/mllike';
import { gleamLanguage }                            from '@exercism/codemirror-lang-gleam';
import { goLanguage }                               from '@codemirror/lang-go';
import { golfScript }                               from 'codemirror-lang-golfscript';
import { groovy }                                   from '@codemirror/legacy-modes/mode/groovy';
import { haskell }                                  from '@codemirror/legacy-modes/mode/haskell';
import { haxe }                                     from '@codemirror/legacy-modes/mode/haxe';
import { j }                                        from 'codemirror-lang-j';
import { janet }                                    from 'codemirror-lang-janet';
import { java }                                     from '@codemirror/lang-java';
import { javascriptLanguage }                       from '@codemirror/lang-javascript';
import { jq }                                       from 'codemirror-lang-jq';
import { julia }                                    from '@codemirror/legacy-modes/mode/julia';
import { k }                                        from 'codemirror-lang-k';
import { lua }                                      from '@codemirror/legacy-modes/mode/lua';
import { nim }                                      from 'nim-codemirror-mode';
import { pascal }                                   from '@codemirror/legacy-modes/mode/pascal';
import { perl }                                     from '@codemirror/legacy-modes/mode/perl';
import { phpLanguage }                              from '@codemirror/lang-php';
import { powerShell }                               from '@codemirror/legacy-modes/mode/powershell';
import { prolog }                                   from 'codemirror-lang-prolog';
import { pythonLanguage }                           from '@codemirror/lang-python';
import { r }                                        from '@codemirror/legacy-modes/mode/r';
import { raku }                                     from './vendor/codemirror-raku';
import { ruby }                                     from '@codemirror/legacy-modes/mode/ruby';
import { rust }                                     from '@codemirror/lang-rust';
import { scheme }                                   from '@codemirror/legacy-modes/mode/scheme';
import { shell }                                    from '@codemirror/legacy-modes/mode/shell';
import { SQLite }                                   from '@codemirror/lang-sql';
import { swift }                                    from '@codemirror/legacy-modes/mode/swift';
import { tcl }                                      from '@codemirror/legacy-modes/mode/tcl';
import { stex }                                     from '@codemirror/legacy-modes/mode/stex';
import { wrenLanguage }                             from '@exercism/codemirror-lang-wren';
import { zig }                                      from 'codemirror-lang-zig';

// Bypass default constructors so we only get highlighters and not extensions.
const elixir     = new LanguageSupport(elixirLanguage);
const gleam      = new LanguageSupport(gleamLanguage);
const go         = new LanguageSupport(goLanguage);
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

// Enable character-wise wrapping whenever possible.
// This was disabled in the upstream due to the old Safari issue (codemirror/dev#524).
const lineWrapping: any = CSS.supports('overflow-wrap', 'anywhere') ? { wordBreak: 'break-all' } : {};

export const extensions : { [key: string]: any } = {
    // Extensions.
    'base': [
        carriageReturn, showUnprintables,
        drawSelection(),
        syntaxHighlighting(defaultHighlightStyle),
        EditorView.theme({
            '.cm-asm-error': { textDecoration: 'underline var(--asm-error)' },
            '.cm-asm-error-tooltip':    asmErrorTooltip,
            '.cm-content':              { fontFamily },
            '.cm-gutters':              { fontFamily },
            '.cm-lineWrapping':         lineWrapping,
            '.cm-tooltip':              { fontFamily },
            '.cm-tooltip-autocomplete': { fontFamily },
        }, { dark: false }),
    ],
    'editor': [
        history(),
        insertChar,
        insertCharState,
        keymap.of([
            // Replace "enter" with a non auto indenting action.
            ...historyKeymap, ...standardKeymap.filter(k => k.key != 'Enter'),
            { key: 'Enter', run: insertNewline },
            { key: 'Tab',   run: insertTab, shift: indentLess },
            { key: 'Mod-/', run: toggleComment },
        ]),
        highlightWhitespace(),
    ],
    'lineNumbers': lineNumbers(),
    'zeroIndexedLineNumbers': lineNumbers(
        {
            formatNumber(num: number) {
                return `${num - 1}`;
            },
        },
    ),
    'bracketMatching': bracketMatching(),
    'vim': vim({ status: true }),

    // Languages.
    // TODO 05ab1e
    // TODO algol-68
    'apl':           StreamLanguage.define(apl),
    // TODO arturo
    'assembly':      assembly(),
    'assembly-wiki': assembly({ byteDumps: false }),
    // TODO awk
    'bash':          StreamLanguage.define(shell),
    // TODO basic
    // TODO befunge
    // TODO berry
    // TODO bqn
    'brainfuck':     brainfuck(),
    'c':             StreamLanguage.define(c),
    'c-sharp':       StreamLanguage.define(csharp),
    'civet':         javascript,
    // TODO cjam
    'clojure':       StreamLanguage.define(clojure),
    'cobol':         StreamLanguage.define(cobol),
    'coconut':       python,
    'coffeescript':  StreamLanguage.define(coffeeScript),
    'common-lisp':   StreamLanguage.define(commonLisp),
    'cpp':           cpp(),
    'crystal':       StreamLanguage.define(crystal),
    'd':             StreamLanguage.define(d),
    'dart':          StreamLanguage.define(dart),
    // TODO egel
    'elixir':        elixir,
    'erlang':        StreamLanguage.define(erlang),
    'f-sharp':       StreamLanguage.define(fSharp),
    'factor':        StreamLanguage.define(factor),
    // TODO fennel
    // TODO fish
    'forth':         StreamLanguage.define({ ...forth, languageData: { commentTokens: { line: '\\' } } }),
    'fortran':       StreamLanguage.define({ ...fortran, languageData: { commentTokens: { line: '!' } } }),
    'gleam':         gleam,
    'go':            go,
    'golfscript':    golfScript(),
    'groovy':        StreamLanguage.define(groovy),
    // TODO harbour
    // TODO hare
    'haskell':       StreamLanguage.define(haskell),
    'haxe':          StreamLanguage.define(haxe),
    // TODO hexagony
    // TODO hush
    // TODO hy
    // TODO iogii
    'j':             j(),
    'janet':         janet(),
    'java':          java(),
    'javascript':    javascript,
    'jq':            jq(),
    'julia':         StreamLanguage.define(julia),
    'k':             k(),
    'kotlin':        StreamLanguage.define(kotlin),
    'lua':           StreamLanguage.define(lua),
    'nim':           StreamLanguage.define({ ...nim( {}, {} ), languageData: { commentTokens: { line: '#' } } }),
    'ocaml':         StreamLanguage.define(oCaml),
    // TODO odin
    'pascal':        StreamLanguage.define(pascal),
    'perl':          StreamLanguage.define(perl),
    'php':           php,
    // TODO picat
    'powershell':    StreamLanguage.define(powerShell),
    'prolog':        prolog(),
    'python':        python,
    // TODO qore
    'r':             StreamLanguage.define(r),
    'racket':        StreamLanguage.define(scheme),
    'raku':          StreamLanguage.define(raku),
    // TODO rebol
    // TODO rexx
    // TODO rockstar
    'ruby':          StreamLanguage.define(ruby),
    'rust':          rust(),
    'scala':         StreamLanguage.define(scala),
    'scheme':        StreamLanguage.define(scheme),
    // TODO sed
    'sql':           sql,
    'squirrel':      StreamLanguage.define(squirrel),
    // TODO stax
    'swift':         StreamLanguage.define(swift),
    'tcl':           StreamLanguage.define(tcl),
    'tex':           StreamLanguage.define(stex),
    // TODO uiua
    // TODO v
    'vala':          StreamLanguage.define(csharp),
    // TODO viml
    // TODO vyxal
    'wren':          wren,
    'zig':           zig(),
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
