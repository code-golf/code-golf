import { $, $$ } from './_util';

const cards   = [...$$('#other-ideas .card')];
const filters = [...$$<HTMLInputElement>('#idea-filters input')];

// Show other ideas whose category is checked, achievement ideas are others.
const filter = () => {
    const checked = new Set(filters.filter(f => f.checked).map(f => f.value));

    for (const card of cards) {
        const category = card.dataset.category ?? 'other';
        card.hidden = !checked.has(category == 'hole' || category == 'lang' ? category : 'other');
    }

    $('#other-ideas-count').innerText = '' + cards.filter(c => !c.hidden).length;
};

for (const f of filters) f.addEventListener('change', filter);
