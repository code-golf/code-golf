function makeLookup(chars: string): Set<number> {
    const set = new Set<number>();
    // for...of iterates over Unicode codepoints (not UTF-16 units)
    for (const c of chars) {
        const codePoint = c.codePointAt(0)!;
        if (codePoint >= 128) {
            set.add(codePoint);
        }
    }
    return set;
}

const allowedStrokesMap: Record<string, Set<number>> = {
    // Method: Open https://tryapl.org/, then run
    // copy([...document.querySelector(".ngn_lb").innerText].filter(c=>c.codePointAt(0)>127).join(""))
    // Note `⎕AV` gives more symbols: special letters like ð (which are legal for identifiers),
    // box drawing characters, and miscellaneous like ¥ and ¶. I see no reason to include these.
    'apl': makeLookup('←×÷⍟⌹○⌈⌊⊥⊤⊣⊢≠≤≥≡≢∨∧⍲⍱↑↓⊂⊃⊆⌷⍋⍒⍳⍸∊⍷∪∩⌿⍀⍪⍴⌽⊖⍉¨⍨⍣∘⍛⍤⍥⍞⎕⍠⌸⌺⌶⍎⍕⋄⍝→⍵⍺∇¯⍬∆⍙'),
    // Method: Open https://www.uiua.org/, then run
    // document.querySelector(".additional-functions").remove()
    // copy([...document.querySelector(".glyph-buttons").innerText].filter(c=>c.codePointAt(0)>127).join("")),
    'uiua': makeLookup('∘◌˙˜⊙⋅⟜⊸⤙⤚◡∩⊃⊓¬±¯⌵√ₑ∿⌊⌈⁅≠≤≥×÷◿ⁿ↧↥∠ℂ⚂ηπτ∞¯←⧻△⇡⊢⊣⇌♭¤⋯⍉⍆⍏⍖⊚◴⊛⧆□⋕≍⊟⊂⊏⊡↯↙↘↻⤸▽⌕⦷∊⨂⊥∧≡⍚⊞⧅⧈⍥⊕⊜◇⌅°⌝⍜⍢⬚⨬⍣⍩⍤'),
    // Method: Open https://github.com/Adriandmen/05AB1E/wiki/Codepage, then run
    // copy([...document.querySelector("table").innerText].filter(c=>c.codePointAt(0)>127).join(""))
    '05ab1e': makeLookup('ǝʒαβγδεζηθвимнтΓΔΘιΣΩ≠∊∍∞₁₂₃₄₅₆Ƶ€Λ‚ƒ„…†‡ˆ‰Š‹ŒĆŽƶĀ‘’“”•–—˜™š›œćžŸā¡¢£¤¥¦§¨©ª«¬λ®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ'),
    // Method Open https://mlochbaum.github.io/BQN/, then run
    // copy([...document.querySelector(".kb").innerText].filter(c=>c.codePointAt(0)>127).join(""))
    'bqn': makeLookup('​​×​÷​⋆​√​⌊​⌈​∧​∨​¬​​≤​​​≥​​≠​≡​≢​⊣​⊢​⥊​∾​≍​⋈​↑​↓​↕​«​»​⌽​⍉​​⍋​⍒​⊏​⊑​⊐​⊒​∊​⍷​⊔​​˙​˜​∘​○​⊸​⟜​⌾​⊘​◶​⎊​⎉​˘​⚇​¨​⌜​⍟​⁼​´​˝​​←​⇐​↩​⋄​​​​​​​​​​⟨​⟩​​​‿​·​•​𝕨​𝕎​𝕩​𝕏​𝕗​𝔽​𝕘​𝔾​𝕤​𝕊​𝕣​¯​π​∞​​​​​'),
    // Method: Open https://github.com/Vyxal/Vyxal/blob/version-3/documentation/codepage.md, then run
    // copy([...document.querySelector("table").innerText].filter(c=>c.codePointAt(0)>127).join(""))
    'vyxal': makeLookup('λƛΛµξ⍾⎋⍟⎊⎄␤⩔Ẅ⊐⎇¿∥∦∺⁜⑴⑵⑶⑷⎂⟒ᛞ▦¨⊞×÷◲⨥⨪ΣΠ⇧⇩∪∩⊍⦰«»ƓɠĠġ⌈⌊⊖⌽£¥↜↝↺↻≜⎀⊢⊣ɦʈᐐᐵᐕ½ƶƵ⁰¹²³⅟※⇄⧖‰≛ℭ℈⦷Ϣ≤≥≠≡•±†⎙γ≓Ͼᴥℳ℗↸⍢ℂ⌹⏚↯⊠⚅æ␣¶★ᑂ∻√⍰◌δ☷σ⎶⊆⍨⎘ꜝ≈≊κ‹›ʀʁɾ▲ṬṪ⤻⤺Ŀ¬∧∨ŁḧᏜᏐ¤⧢①②③④⑤⑥⑦⑧Þ∆ø„“”'),
};

/**
 * Return a set of unicode codepoints that count as 1 stroke rather than the
 * number of bytes in their UTF-8 representations.
 * Return `undefined` if the language should never use 'strokes' scoring.
 */
export function getAllowedStrokes(lang: string): Set<number> | undefined {
    return allowedStrokesMap[lang];
}
