import { ASMStateField, ASMLanguageData, byteDumper, ASMColorFacet, SectionColors, ASMFlush }    from "./compilerPlugin.js";
import { errorMarker, errorTooltipper }    from "./errorPlugin.js";
import { parser }                          from "./parser.js";
import { debugPlugin }                     from "./debugPlugin.js";
import { ShellcodePlugin, ShellcodeField } from "./shellcodePlugin.js";
import { ctxTracker }                   from "./tokenizer.js";
import { LRLanguage, LanguageSupport }  from '@codemirror/language';
import { styleTags, tags }              from '@lezer/highlight';

const assemblyLang = LRLanguage.define({
    parser: parser.configure({
        props: [
            styleTags({
                Opcode: tags.operatorKeyword,
                IOpcode: tags.operatorKeyword,
                RelOpcode: tags.operatorKeyword,
                IRelOpcode: tags.operatorKeyword,
                Prefix: tags.operatorKeyword,
                Register: tags.className,
                Directive: tags.meta,
                Comment: tags.lineComment,
                SymbolName: tags.definition(tags.labelName),
                Immediate: tags.literal,
                IImmediate: tags.literal,
                Memory: tags.regexp,
                IMemory: tags.regexp,
                Relative: tags.regexp,
                Expression: tags.literal,
                FullString: tags.string,
                VEXRound: tags.modifier,
                VEXMask: tags.modifier,
                Offset: tags.emphasis,
                Ptr: tags.emphasis,
                SpecialWord: tags.annotation
            })
        ]
    })
});

/** Create a CodeMirror extension utilizing DefAssembler
 * @param {Object} assemblyState The assembly state. Must be compatible with DefAsm's AssemblyState plus a Range member that must be a constructor for ranges (like DefAsm's Range).
 * @param {Object} config
 * @param {boolean} config.byteDumps Whether to display the results of the assembly on the side of the editor
 * @param {boolean} config.debug Whether to enable toggling debug mode when pressing F3
 * @param {boolean} config.errorMarking Whether to draw a red underline beneath segments of code that cause errors
 * @param {boolean} config.errorTooltips Whether to show a tooltip on these segments explaining the cause of the error
 * @param {boolean} config.highlighting Whether to enable syntax highlighting using a [`LanguageSupport`](https://codemirror.net/6/docs/ref/#language.LanguageSupport) object
*/
export function assemblyIde(assemblyState, {
    byteDumps      = true,
    debug          = false,
    errorMarking   = true,
    errorTooltips  = true,
} = {})
{
    const plugins = [
        ASMStateField.init(state => {
            const asm = assemblyState;
            asm.compile(state.sliceDoc());
            return asm;
        }),
        ASMLanguageData
    ];
    if(byteDumps)     plugins.push(SectionColors, byteDumper);
    if(debug)         plugins.push(debugPlugin);
    if(errorMarking)  plugins.push(errorMarker);
    if(errorTooltips) plugins.push(errorTooltipper);

    return plugins;
}

export function assemblyLanguage(bitness        = 64,
    defaultSyntax  = {
        intel: false,
        prefix: true
    })
{
    return new LanguageSupport(assemblyLang.configure({
        contextTracker: ctxTracker(defaultSyntax, bitness)
    }), []);
}

export { ASMStateField, ASMColorFacet, ASMFlush, ShellcodePlugin, ShellcodeField };
