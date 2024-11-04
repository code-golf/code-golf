import { EditorState, EditorView, extensions } from './_codemirror';
import { $$ }                                  from './_util';
import 'mathjax/es5/tex-chtml.js';

export function highlightCodeBlocks(selector: string){
    const baseExtensions = [...extensions.base, EditorState.readOnly.of(true)];

    for (const code of $$(selector)) {
        const lang = {
            'asm' : 'assembly',
            'c++' : 'cpp',
            'cs'  : 'c-sharp',
            'fs'  : 'f-sharp',
            'ijs' : 'j',
        }[code.className.slice('language-'.length).toLowerCase()] || '';

        // Clear the existing code and replace with a read-only editor.
        const pre = code.parentElement!;
        pre.innerHTML = '';
        new EditorView({
            parent: pre,
            state:  EditorState.create({
                doc: code.innerText.trimEnd(),
                extensions: [
                    baseExtensions,
                    extensions[`${lang}-wiki`] ?? extensions[lang] ?? [],
                ],
            }),
        });
    }
}
