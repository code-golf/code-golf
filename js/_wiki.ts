import { EditorState, EditorView, extensions } from './_codemirror';
import { $$ }                                  from './_util';
import 'mathjax/es5/tex-chtml.js';

export function highlightCodeBlocks(selector: string){
    const baseExtensions = [...extensions.base, EditorState.readOnly.of(true)];

    for (const code of $$(selector)) {
        let lang = code.className.slice('language-'.length).toLowerCase();

        // Should be handled instead by MathJax.
        if (lang == 'math') continue;

        // FIXME Is there a better way to do this mapping?
        if (lang == 'asm') lang = 'assembly';
        if (lang == 'c++') lang = 'cpp';
        if (lang == 'cs' ) lang = 'c-sharp';
        if (lang == 'fs' ) lang = 'f-sharp';
        if (lang == 'ijs') lang = 'j';

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
