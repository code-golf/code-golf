import {
    ComponentItem, ComponentItemConfig, ContentItem, GoldenLayout,
    RowOrColumn, LayoutConfig, ResolvedRootItemConfig,
    ResolvedLayoutConfig, DragSource, LayoutManager, ComponentContainer,
} from 'golden-layout';
import { EditorView }                     from './_codemirror';
import diffTable                          from './_diff';
import { $, $$, byteLen, charLen, comma } from './_util';
import {
    init, langs, getLang, hole, getAutoSaveKey, setSolution, getSolution,
    setCode, refreshScores, getHideDeleteBtn, submit, SubmitResponse,
    updateRestoreLinkVisibility, getSavedInDB, setCodeForLangAndSolution,
    populateScores, getCurrentSolutionCode, initDeleteBtn, initCopyJSONBtn,
    getScorings,
} from './_hole-common';

const poolDragSources: {[key: string]: DragSource} = {};
const poolElements: {[key: string]: HTMLElement} = {};

/**
 * Is mobile mode activated? Start at false as default since Golden Layout
 * uses desktop as default. Change to true and apply changes if width is less
 * than or equal to 768px (it seems to be a common breakpoint idk).
 *
 * Changes on mobile mode:
 * - golden layout reflowed to columns-only
 * - full-page scrolling is enabled
 * - dragging is disabled (incompatible with being able to scroll)
 * - maximized windows take the full screen
 *
 * TODO: respect "Request desktop site" from mobile browsers to force
 * isMobile = false. Or otherwise configuration option.
 */
let isMobile = false;
let applyingDefault = false;

let subRes: SubmitResponse | null = null;
const readonlyOutputs: {[key: string]: HTMLElement | undefined} = {};

let editor: EditorView | null = null;

init(true, setSolution, setCodeForLangAndSolution, updateReadonlyPanels, () => editor);

// Handle showing/hiding alerts
for (const alertCloseBtn of $$('.main_close')) {
    const alert = alertCloseBtn.parentNode as HTMLDivElement;
    alertCloseBtn.addEventListener('click', () => {
        const child = (alert.querySelector('svg') as any).cloneNode(true);
        $('#alert-pool').appendChild(child);
        alert.classList.add('hide');
        child.addEventListener('click', () => {
            child.parentNode.removeChild(child);
            alert.classList.remove('hide');
        });
    });
}

// Handle showing/hiding lang picker
// can't be done in CSS because the picker is one parent up
const langToggle = $<HTMLDetailsElement>('#hole-lang details');
langToggle.addEventListener('toggle', () => {
    $('#picker').classList.toggle('hide', !langToggle.open);
});

const goldenContainer = $('#golden-container');

/**
 * Actual Golden Layout docs are at
 *  https://golden-layout.github.io/golden-layout
 * golden-layout.com is for the old GL.
 */
const layout = new GoldenLayout(goldenContainer);
layout.resizeWithContainerAutomatically = true;

function updateReadonlyPanel(name: string) {
    if (!subRes) return;
    const output = readonlyOutputs[name];
    if (!output) return;
    switch (name) {
    case 'err':
        output.innerHTML = subRes.Err.replace(/\n/g,'<br>');
        break;
    case 'out':
        output.innerText = subRes.Out;
        break;
    case 'exp':
        output.innerText = subRes.Exp;
        break;
    case 'arg':
        // Hide arguments unless we have some.
        output.replaceChildren(
            ...subRes.Argv.map(a => <span>{a}</span>),
        );
        break;
    case 'diff':
        const diff = diffTable(hole, subRes.Exp, subRes.Out, subRes.Argv);
        output.replaceChildren(diff);
    }
}

function updateReadonlyPanels(data: SubmitResponse) {
    subRes = data;
    for (const name in readonlyOutputs) {
        updateReadonlyPanel(name);
    }
}

for (const i of [0,1,2,3,4]) {
    const name = ['exp', 'out', 'err', 'arg', 'diff'][i];
    const title = ['Expected', 'Output', 'Errors', 'Arguments', 'Diff'][i];
    layout.registerComponentFactoryFunction(name, container => {
        container.setTitle(title);
        autoFocus(container);
        container.element.id = name;
        container.element.classList.add('readonly-output');
        readonlyOutputs[name] = container.element;
        updateReadonlyPanel(name);
    });
}

function makeEditor(parent: HTMLDivElement) {
    editor = new EditorView({
        dispatch: tr => {
            if (!editor) return;
            const result = editor.update([tr]) as unknown;

            const code = tr.state.doc.toString();
            const scorings: {total: {byte?: number, char?: number}, selection?: {byte?: number, char?: number}} = getScorings(tr, editor);
            const scoringKeys = ['byte', 'char'] as const;

            function formatScore(scoring) {
                return scoringKeys
                    .filter(s => s in scoring)
                    .map(s => `${comma(scoring[s])} ${s}${scoring[s] != 1 ? 's' : ''}`)
                    .join(', ');
            }

            if (scorings.selection) {
                $('#strokes').innerText = `${formatScore(scorings.total)} (${formatScore(scorings.selection)} selected)`;
            } else {
                $('#strokes').innerText = formatScore(scorings.total);
            }

            // Avoid future conflicts by only storing code locally that's
            // different from the server's copy.
            const serverCode = getCurrentSolutionCode();

            const key = getAutoSaveKey(getLang(), getSolution());
            if (code && (code !== serverCode || !getSavedInDB()) && code !== langs[getLang()].example)
                localStorage.setItem(key, code);
            else
                localStorage.removeItem(key);

            updateRestoreLinkVisibility(editor);

            return result;
        },
        parent: parent,
    });

    editor.contentDOM.setAttribute('data-gramm', 'false');  // Disable Grammarly.
}

function autoFocus(container: ComponentContainer) {
    container.element.addEventListener('focusin', () => container.focus());
    container.element.addEventListener('click', () => container.focus());
}

layout.registerComponentFactoryFunction('code', async container => {
    container.setTitle('Code');
    autoFocus(container);

    const header = (<header>
        <div id="strokes">0 bytes, 0 chars</div>
        <a class="hide" href="" id="restoreLink">Restore solution</a>
    </header>) as HTMLElement;
    const editorDiv = <div id="editor"></div> as HTMLDivElement;

    makeEditor(editorDiv);

    header.append($<HTMLTemplateElement>('#template-run').content.cloneNode(true));

    container.element.id = 'editor-section';
    container.element.append(editorDiv, header);

    await afterDOM();

    $('#restoreLink').onclick = (e: MouseEvent) => {
        e.preventDefault();
        setCode(getCurrentSolutionCode(), editor);
    };

    // Wire submit to clicking a button and a keyboard shortcut.
    $('#runBtn').onclick = () => submit(editor, updateReadonlyPanels);

    const deleteBtn = $('#deleteBtn');
    if (deleteBtn) {
        initDeleteBtn(deleteBtn, langs);
        deleteBtn.classList.toggle('hide', getHideDeleteBtn());
    }

    setCodeForLangAndSolution(editor);
});

async function afterDOM() {}

function delinkRankingsView() {
    $$('#rankingsView a').forEach(a => a.onclick = e => {
        e.preventDefault();

        $$<HTMLAnchorElement>('#rankingsView a').forEach(a => a.href = '');
        a.removeAttribute('href');

        document.cookie =
            `rankings-view=${a.innerText.toLowerCase()};SameSite=Lax;Secure`;

        refreshScores(editor);
    });
}

layout.registerComponentFactoryFunction('scoreboard', async container => {
    container.setTitle('Scoreboard');
    autoFocus(container);
    container.element.append(
        $<HTMLTemplateElement>('#template-scoreboard').content.cloneNode(true),
    );
    container.element.id = 'scoreboard-section';
    await afterDOM();
    populateScores(editor);
    delinkRankingsView();
});

layout.registerComponentFactoryFunction('details', container => {
    container.setTitle('Details');
    autoFocus(container);
    const details = $<HTMLTemplateElement>('#template-details').content.cloneNode(true) as HTMLDetailsElement;
    container.element.append(details);
    container.element.id = 'details-content';
    initCopyJSONBtn(container.element.querySelector('#copy') as HTMLElement);
});

function plainComponent(componentType: string): ComponentItemConfig {
    return {
        type: 'component',
        componentType: componentType,
        reorderEnabled: !isMobile,
    };
}

const defaultLayout: LayoutConfig = {
    settings: {
        showPopoutIcon: false,
    },
    dimensions: {
        headerHeight: 28,
    },
    root: {
        type: 'column',
        content: [
            {
                type: 'row',
                content: [
                    {
                        ...plainComponent('code'),
                        width: 75,
                    },
                    {
                        ...plainComponent('scoreboard'),
                        width: 25,
                    },
                ],
            }, {
                type: 'row',
                content: [
                    {
                        type: 'stack',
                        content: [
                            plainComponent('arg'),
                            plainComponent('exp'),
                        ],
                    }, {
                        type: 'stack',
                        content: [
                            plainComponent('out'),
                            plainComponent('err'),
                            plainComponent('diff'),
                        ],
                    },
                ],
            },
        ],
    },
};

async function applyDefaultLayout() {
    applyingDefault = true;
    toggleMobile(false);
    Object.keys(poolElements).map(removePoolItem);
    addPoolItem('details', 'Details');
    layout.loadLayout(defaultLayout);
    await afterDOM();
    checkMobile();
    applyingDefault = false;
}

applyDefaultLayout();

/**
 * Try to add after selected item, with sensible defaults
 */
function addItemFromPool(componentName: string) {
    (window as any).layout = layout;
    layout.addItemAtLocation(
        plainComponent(componentName),
        LayoutManager.afterFocusedItemIfPossibleLocationSelectors,
    );
}

/**
 * Add the first element from the pool to the root column, or create a new
 * column containing the root and the first pool element if non exist.
 */
function addRow() {
    if (!layout.rootItem) return;
    const newComponentName = Object.keys(poolElements)[0];
    if (!newComponentName) return;
    const newConfig = plainComponent(newComponentName);
    if (layout.rootItem.type === 'column') {
        // Add to existing column
        (layout.rootItem as RowOrColumn).addItem(newConfig);
    }
    else {
        // Create new column
        const oldParent = layout.rootItem;
        const newParent = (layout as any).createContentItem({
            type: 'column',
            content: [],
        });
        const oldParentParent = oldParent.parent!;
        // removeChild(_, true): don't remove the node entirely, just remove
        // it from the current tree before re-inserting
        oldParentParent.removeChild(oldParent, true);
        oldParentParent.addChild(newParent);
        newParent.addChild(oldParent);
        newParent.addItem(newConfig);
    }
    (layout as any).getAllContentItems().find(
        (item: ComponentItem) => item.componentType === newComponentName,
    )?.focus();
}

$('#add-row').addEventListener('click', addRow);

$('#revert-layout').addEventListener('click', applyDefaultLayout);

$('#make-wide').addEventListener('click',
    () => document.documentElement.classList.toggle('full-width', true),
);

$('#make-narrow').addEventListener('click',
    () => document.documentElement.classList.toggle('full-width', false),
);

function addPoolItem(componentType: string, title: string) {
    poolElements[componentType]?.remove();
    const el = (<span class="btn">{title}</span>);
    $('#pool').appendChild(el);
    poolDragSources[componentType] = layout.newDragSource(el, componentType);
    poolElements[componentType] = el;
    checkShowAddRow();
    el.addEventListener('click', () => addItemFromPool(componentType));
}

// Add an item to the tab pool when a component gets destroyed
layout.addEventListener('itemDestroyed', e => {
    if (applyingDefault) return;
    const _target = e.target as ContentItem;
    if (_target.isComponent) {
        const target = _target as ComponentItem;
        addPoolItem(target.componentType as string, target.title);
    }
    checkShowAddRow();
});

function removePoolItem(componentType: string) {
    if (!poolElements[componentType]) return;
    poolElements[componentType].remove();
    delete poolElements[componentType];
    checkShowAddRow();
    if (!isMobile) removeDragSource(componentType);
}

async function checkShowAddRow() {
    // Await to ensure that rootItem === undefined after removing last item
    await afterDOM();
    $('#add-row').classList.toggle(
        'hide',
        Object.keys(poolElements).length === 0
            || layout.rootItem === undefined,
    );
}

function removeDragSource(componentType: string) {
    layout.removeDragSource(poolDragSources[componentType]);
    delete poolDragSources[componentType];
}

// Remove an item from the tab pool when it gets added
layout.addEventListener('itemCreated', e => {
    if (applyingDefault) return;
    const target = e.target as ContentItem;
    if (target.isComponent) {
        removePoolItem((target as ComponentItem).componentType as string);
    }
});


/**
 * There's a bug with the dragging from layout.newDragSource where dragging up
 * from the tab pool causes a .lm_dragProxy to appear, but it doesn't get
 * removed due to an error "Ground node can only have a single child." Rather
 * than fix the bug, just remove all .lm_dragProxy elements after mouseups that
 * follow a state change.
 *
 * The error message still gets logged in console
 */
layout.addEventListener('stateChanged', () => {
    document.addEventListener('mouseup', removeDragProxies);
    document.addEventListener('touchend', removeDragProxies);
    document.documentElement.classList.toggle('has_lm_maximised', !!$('.lm_maximised'));
});

function removeDragProxies() {
    $$('.lm_dragProxy').forEach(e => e.remove());
    document.removeEventListener('mouseup', removeDragProxies);
    document.removeEventListener('touchend', removeDragProxies);
}

/**
 * LayoutConfig has a bunch of optional properties, while ResolvedLayoutConfig
 * marks everything as readonly for no reason. We converted ResolvedLayoutConfig
 * to a superset of LayoutConfig by making everything mutable.
 */
type DeepMutable<T> = { -readonly [key in keyof T]: DeepMutable<T[key]> };

/**
 * Mutate the given item recursively to:
 * - change reorderEnabled (false if isMobile, otherwise true)
 * - change rows to columns (if isMobile, otherwise no change)
 *
 * I don't know what it's necessary to change reorderEnabled on a per-item
 * basis. Should be able to just do currLayout.settings.reorderEnabled = ...,
 * but that is not respected at all, even for new items.
 */
function mutateDeep(item: DeepMutable<ResolvedRootItemConfig>, isMobile: boolean) {
    if (isMobile && item.type === 'row') {
        (item as any).type = 'column';
    }
    (item as any).reorderEnabled = !isMobile;
    if (item.content.length > 0) {
        item.content.forEach(child => mutateDeep(child, isMobile));
    }
}

function toggleMobile(_isMobile: boolean) {
    isMobile = _isMobile;
    // This could be a CSS media query, but I'm keeping generality in case of
    // other config options ("request desktop site", button config, etc.)
    document.documentElement.classList.toggle('mobile', isMobile);
    const currLayout = layout.saveLayout() as DeepMutable<ResolvedLayoutConfig>;
    if (currLayout.root) {
        mutateDeep(currLayout.root, isMobile);
        layout.loadLayout(currLayout as any as LayoutConfig);
    }
    if (isMobile) {
        for (const componentType in poolDragSources)
            removeDragSource(componentType);
    }
    else {
        for (const componentType in poolElements)
            poolDragSources[componentType] = layout.newDragSource(poolElements[componentType], componentType);
    }
    updateMobileContainerHeight();
}

function checkMobile() {
    if ((window.innerWidth < 768) !== isMobile) {
        toggleMobile(!isMobile);
    }
}

window.addEventListener('resize', checkMobile);

/**
 * Golden Layout has handlers for both "touchstart" and "click," which is a
 * problem because a touch on mobile triggers both events (example symptom:
 * tapping "close" button closes two tabs instead of one).
 *
 * Duplicate handlers are present on:
 * - header maximize/close buttons
 * - tab "close" button
 * - tab itself (doesn't matter because selection is idempotent)
 * - header bar (doesn't matter because we don't use it)
 *
 * We work around this by going into GL internals and disabling the touchstart
 * callbacks. This is not supported behavior, but it works.
 */
function deepCancelTouchStart(item: any) {
    if (!item) return;
    if (item.type === 'stack') {
        item._header._closeButton.onTouchStart = () => {};
        item._header._maximiseButton.onTouchStart = () => {};
    }
    else if (item.type === 'component') {
        item._tab.onCloseTouchStart = () => {};
    }
    item._contentItems?.forEach((child: any) => deepCancelTouchStart(child));
}

deepCancelTouchStart(layout.rootItem);

layout.addEventListener('stateChanged', () => {
    deepCancelTouchStart(layout.rootItem);
    updateMobileContainerHeight();
});

function rowCount(item: ContentItem | undefined): number {
    if (!item) return 0;
    if (item.type === 'row')
        return Math.max(...item.contentItems.map(rowCount));
    else if (item.type === 'column')
        return item.contentItems.map(rowCount).reduce((a, b) => a + b);
    else
        return 1;
}

function updateMobileContainerHeight() {
    goldenContainer.style.height =
        isMobile ? rowCount(layout.rootItem) * 25 + 'rem' : '';
}
