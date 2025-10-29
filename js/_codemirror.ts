import { EditorState, Extension }                         from '@codemirror/state';
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

import { AssemblyState as DefAsmState }             from '@defasm/core';
import { AssemblyState as ExtAsmState }             from './assembly/extasm';

// Languages.
import { apl }                                      from '@codemirror/legacy-modes/mode/apl';
import { assemblyLanguage, assemblyIde }            from './assembly/codemirror';
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
import { knight }                                   from 'codemirror-lang-knight';
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
import { createAssembler, PreprocessedAssembler, ToolchainAssembler } from './assembly/assemble';

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

function x86Assembly(byteDumps: boolean = true) {
    const state = new DefAsmState();
    return [assemblyLanguage({bitness: state.bitness, defaultSyntax: state.defaultSyntax}), assemblyIde(state, {byteDumps})];
}

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
    'arm64':      assemblyLanguage(),
    'assembly':      x86Assembly(),
    'assembly-wiki': x86Assembly(false),
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
    'knight':        knight(),
    'kotlin':        StreamLanguage.define(kotlin),
    'lua':           StreamLanguage.define(lua),
    'luau':          StreamLanguage.define(lua),
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
    'risc-v':      assemblyLanguage(),

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
    // TODO umka
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

import aspp from '../wasm/aspp.wasm';
import risc_v_as from '../wasm/riscv64-linux-gnu-as.wasm';
import risc_v_ld from '../wasm/riscv64-linux-gnu-ld.wasm';
import risc_v_objdump from '../wasm/riscv64-linux-gnu-objdump.wasm';
import arm64_as from '../wasm/aarch64-linux-gnu-as.wasm';
import arm64_ld from '../wasm/aarch64-linux-gnu-ld.wasm';
import arm64_objdump from '../wasm/aarch64-linux-gnu-objdump.wasm';
import { wasmCompileUrl } from './assembly/wasm-program';

const linkerScript = `
ENTRY(_start)

SECTIONS
{
	. = 0x400000;

	.text : ALIGN(0x1000) {
		*(.text)
	} :text

    .debug          0 : { *(.debug) }
    .line           0 : { *(.line) }
    .debug_srcinfo  0 : { *(.debug_srcinfo) }
    .debug_sfnames  0 : { *(.debug_sfnames) }
    .debug_aranges  0 : { *(.debug_aranges) }
    .debug_pubnames 0 : { *(.debug_pubnames) }
    .debug_info     0 : { *(.debug_info .gnu.linkonce.wi.*) }
    .debug_abbrev   0 : { *(.debug_abbrev) }
    .debug_line     0 : { *(.debug_line .debug_line.* .debug_line_end) }
    .debug_frame    0 : { *(.debug_frame) }
    .debug_str      0 : { *(.debug_str) }
    .debug_loc      0 : { *(.debug_loc) }
    .debug_macinfo  0 : { *(.debug_macinfo) }
    .debug_weaknames 0 : { *(.debug_weaknames) }
    .debug_funcnames 0 : { *(.debug_funcnames) }
    .debug_typenames 0 : { *(.debug_typenames) }
    .debug_varnames  0 : { *(.debug_varnames) }
    .debug_pubtypes 0 : { *(.debug_pubtypes) }
    .debug_ranges   0 : { *(.debug_ranges) }
    .debug_addr     0 : { *(.debug_addr) }
    .debug_line_str 0 : { *(.debug_line_str) }
    .debug_loclists 0 : { *(.debug_loclists) }
    .debug_macro    0 : { *(.debug_macro) }
    .debug_names    0 : { *(.debug_names) }
    .debug_rnglists 0 : { *(.debug_rnglists) }
    .debug_str_offsets 0 : { *(.debug_str_offsets) }
    .debug_sup      0 : { *(.debug_sup) }

	.data : ALIGN(0x1000) {
		*(.data)
		*(.rodata*)
		*(.*)
	} :data

	.bss : ALIGN(0x1000) {
		*(.bss)
	} :bss
}

PHDRS
{
	text PT_LOAD FLAGS(7) AT(0x400000);
	data PT_LOAD FLAGS(7);
	bss PT_LOAD FLAGS(7);
};
`

function riscvAsArgs(src: string) {
    let arch;
    let archLine = false;
    if(src.startsWith("#32")) {
        archLine = true;
        arch = "rv32";
    } else {
        if(src.startsWith("#64")) {
            archLine = true;
        }
        arch = "rv64";
    }
    arch += "gmafdqcbvh_zic64b_ziccamoa_ziccif_zicclsm_ziccrse_zicbom_zicbop_zicboz_zicond_zicntr_zicsr_zifencei_zihintntl_zihintpause_zihpm_zimop_zicfiss_zicfilp_zmmul_za64rs_za128rs_zaamo_zabha_zacas_zalrsc_zawrs_zfbfmin_zfa_zfh_zfhmin_zbb_zba_zbc_zbs_zbkb_zbkc_zbkx_zk_zkn_zknd_zkne_zknh_zkr_zks_zksed_zksh_zkt_zve32x_zve32f_zve64x_zve64f_zve64d_zvbb_zvbc_zvfbfmin_zvfbfwma_zvfh_zvfhmin_zvkb_zvkg_zvkn_zvkng_zvknc_zvkned_zvknha_zvknhb_zvksed_zvksh_zvks_zvksg_zvksc_zvkt_zvl32b_zvl64b_zvl128b_zvl256b_zvl512b_zvl1024b_zvl2048b_zvl4096b_zvl8192b_zvl16384b_zvl32768b_zvl65536b_ztso_zca_zcb_zcmop_xtheadba_xtheadbb_xtheadbs_xtheadcmo_xtheadcondmov_xtheadfmemidx_xtheadfmv_xtheadmac_xtheadmemidx_xtheadmempair_xtheadsync_xventanacondops";
    let preferM = false;
    if(archLine) {
        for(let i = 1; i < src.length; ++i) {
            const c = src[i];
            if(c == '\n') {
                break;
            }
            if(c == 'M') {
                preferM = true;
            }
        }
    }
    if(preferM) {
        arch += "_zcmp_zcmt";
    } else {
        arch += "_zcd";
    }
    return ["-march", arch, "-mno-arch"];
}

function riscvLdArgs(src: string) {
    const machine = src.startsWith("#32") ? "elf32lriscv" : "elf64lriscv";
    return ["-m", machine];
}

function arm64AsArgs(src: string) {
    return ["-march=armv9.5-a+fp+bf16+crypto+crc+f32mm+f64mm+fp16+fp16fml+memtag+rng+sb+simd+sme+sme-f64f64+sme-i16i64+sve+sve2"];
}

function arm64LdArgs(src: string) {
    return [];
}

function assemblyIdeFactory(as: string, ld: string, objdump: string, asArgs: (src: string) => string[], ldArgs: (src: string) => string[]) {
    const ldArgs2 = (src: string) => {const args = ldArgs(src); args.push("--no-warn-rwx"); return args;};
    return async () => assemblyIde(new ExtAsmState(await PreprocessedAssembler.create(wasmCompileUrl(aspp),
        ToolchainAssembler.create(wasmCompileUrl(as), wasmCompileUrl(ld), wasmCompileUrl(objdump)).then(x =>
            x.withLinkerScript(linkerScript).withAsArgs(asArgs).withLdArgs(ldArgs2)
        )
    )));
}

export const asyncExtensions : { [key: string]: () => Promise<Extension> } = {
    'arm64': assemblyIdeFactory(arm64_as, arm64_ld, arm64_objdump, arm64AsArgs, arm64LdArgs),
    'risc-v': assemblyIdeFactory(risc_v_as, risc_v_ld, risc_v_objdump, riscvAsArgs, riscvLdArgs),
};
