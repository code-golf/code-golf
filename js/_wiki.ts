import { EditorState, EditorView, extensions } from './_codemirror';
import { $$ }                                  from './_util';

export function highlightCodeBlocks(selector: string){
    const baseExtensions = [...extensions.base, EditorState.readOnly.of(true)];

    for (const code of $$(selector)) {
        let lang = code.className.slice('language-'.length).toLowerCase();

        // FIXME Is there a better way to do this mapping?
        if (lang == 'asm') lang = 'assembly';
        if (lang == 'c++') lang = 'cpp';
        if (lang == 'fs' ) lang = 'f-sharp';
        if (lang == 'ijs') lang = 'j';

        // Skip Assembly for now as the annoations break the layout.
        if (lang == 'assembly') continue;

        // Clear the existing code and replace with a read-only editor.
        const pre = code.parentElement!;
        pre.innerHTML = '';
        new EditorView({
            parent: pre,
            state:  EditorState.create({
                doc: code.innerText.trim(),
                extensions: [
                    baseExtensions,
                    extensions[lang as keyof typeof extensions] ?? [],
                ],
            }),
        });
    }
}
