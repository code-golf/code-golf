import {
    ComponentItem, ComponentItemConfig, ContentItem, GoldenLayout,
    RowOrColumn, LayoutConfig, ResolvedRootItemConfig,
    DragSource, LayoutManager, ComponentContainer, ResolvedLayoutConfig,
    RootItemConfig,
} from 'golden-layout';
import { EditorView }   from './_codemirror';
import diffView         from './_diff';
import { $, $$, comma, throttle } from './_util';
import {
    init, langs, hole, setSolution,
    setCode, refreshScores, getHideDeleteBtn, submit, ReadonlyPanelsData,
    updateRestoreLinkVisibility, setCodeForLangAndSolution,
    populateScores, getCurrentSolutionCode, initDeleteBtn, initCopyButtons,
    getScorings,
    updateLocalStorage,
    getLang,
    setState,
    ctrlEnter,
    getLastSubmittedCode,
} from './_hole-common';
import { highlightCodeBlocks } from './_wiki';
import UnprintableElement from './_unprintable';

const poolDragSources: {[key: string]: DragSource} = {};
const poolElements: {[key: string]: HTMLElement} = {};
let isWide = false;

const isSponsor = JSON.parse($('#sponsor').innerText);
const hasNotes  = JSON.parse($('#has-notes').innerText);

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

let subRes: ReadonlyPanelsData | null = null;
let langWikiContent = '';
let holeLangNotesContent = '';
const readonlyOutputs: {[key: string]: HTMLElement | undefined} = {};

let editor: EditorView | null = null;
let holeLangNotesEditor: EditorView | null = null;

let substitutions: {pattern: RegExp, replacement: string}[] = [];

init(true, setSolution, setCodeForLangAndSolution, updateReadonlyPanels, () => editor);

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
        output.replaceChildren(UnprintableElement.escape(subRes.Out));
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
        const ignoreCase = JSON.parse($('#case-fold').innerText);
        const diff = diffView(hole, subRes.Exp, subRes.Out, subRes.Argv, ignoreCase, subRes.OutputDelimiter, subRes.MultisetItemDelimiter);
        output.replaceChildren(diff);
    }
}

function updateWikiContent() {
    if ($('#langWiki')) {
        $('#langWiki').innerHTML = `<article>${langWikiContent}</article>`;
        highlightCodeBlocks('#langWiki pre > code');
    }
}

function updateHoleLangNotesContent() {
    const input = $<HTMLInputElement>('#notes-substitutions');
    if (input) input.value = localStorage.getItem(`${getLang()}-substitutions`) ?? '';
    if (holeLangNotesEditor) setState(holeLangNotesContent, holeLangNotesEditor);
}

function updateReadonlyPanels(data: ReadonlyPanelsData | {langWiki: string}| {holeLangNotes: string}) {
    if ('langWiki' in data) {
        langWikiContent = data.langWiki;
        updateWikiContent();
    }
    else if ('holeLangNotes' in data) {
        holeLangNotesContent = data.holeLangNotes;
        updateHoleLangNotesContent();
    }
    else {
        subRes = data;
        for (const name in readonlyOutputs) {
            updateReadonlyPanel(name);
        }
    }
}

for (const name of ['exp', 'out', 'err', 'arg', 'diff']) {
    layout.registerComponentFactoryFunction(name, container => {
        container.setTitle(getTitle(name));
        autoFocus(container);
        container.element.id = name;
        container.element.classList.add('readonly-output');
        readonlyOutputs[name] = container.element;
        updateReadonlyPanel(name);
    });
}

layout.registerComponentFactoryFunction('langWiki', async container => {
    container.setTitle(getTitle('langWiki'));
    autoFocus(container);
    container.element.id = 'langWiki';
    await afterDOM();
    updateWikiContent();
});

function makeEditor(parent: HTMLDivElement) {
    editor = new EditorView({
        dispatch: tr => {
            if (!editor) return;
            const result = editor.update([tr]) as unknown;

            const code = tr.state.doc.toString();
            const scorings: {total: {byte?: number, char?: number}, selection?: {byte?: number, char?: number}} = getScorings(tr, editor);
            const scoringKeys = ['byte', 'char'] as const;

            $('main')?.classList.toggle('lastSubmittedCode', code === getLastSubmittedCode());

            function formatScore(scoring: any) {
                return scoringKeys
                    .filter(s => s in scoring)
                    .map(s => `${comma(scoring[s])} ${s}${scoring[s] != 1 ? 's' : ''}`)
                    .join(', ');
            }

            $('#strokes').innerText = (() => {
                let innertext = scorings.selection
                    ? `${formatScore(scorings.total)} (${formatScore(scorings.selection)} selected)`
                    : formatScore(scorings.total);
                if (scorings.selection?.char === 1) {
                    const sel = tr.state.sliceDoc(tr.state.selection.main.from, tr.state.selection.main.to);
                    const code = sel.codePointAt(0);
                    if (code !== undefined) {
                        const hex = code.toString(16).toUpperCase().padStart(4, '0');
                        const bin = code.toString(2).padStart(8, '0');
                        innertext += ` ${sel} - ${code} - 0x${hex} - ${bin}`;
                    }
                }
                return innertext;
            })();


            updateLocalStorage(code);
            updateRestoreLinkVisibility(editor);

            return result;
        },
        parent,
    });

    editor.contentDOM.setAttribute('data-gramm', 'false');  // Disable Grammarly.
}

function makeHoleLangNotesEditor(parent: HTMLDivElement) {
    holeLangNotesEditor = new EditorView({
        dispatch: tr => {
            if (!holeLangNotesEditor) return;
            const result = holeLangNotesEditor.update([tr]) as unknown;
            const content = tr.state.doc.toString();
            $<HTMLButtonElement>('#upsert-notes-btn').disabled = content === holeLangNotesContent || (!!content && !isSponsor);
            return result;
        },
        parent,
    });

    holeLangNotesEditor.contentDOM.setAttribute('data-gramm', 'false');  // Disable Grammarly.
}

function autoFocus(container: ComponentContainer) {
    container.element.addEventListener('focusin', () => container.focus());
    container.element.addEventListener('click', () => container.focus());
}

layout.registerComponentFactoryFunction('code', async container => {
    container.setTitle(getTitle('code'));
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
    const closuredSubmit = () => submit(editor, updateReadonlyPanels);
    $('#runBtn').onclick = closuredSubmit;
    $('#editor').onkeydown = ctrlEnter(closuredSubmit);

    const deleteBtn = $('#deleteBtn');
    if (deleteBtn) {
        initDeleteBtn(deleteBtn, langs);
        deleteBtn.classList.toggle('hide', getHideDeleteBtn());
    }

    setCodeForLangAndSolution(editor);
});

async function upsertNotes() {
    $<HTMLButtonElement>('#upsert-notes-btn').disabled = true;
    const content = holeLangNotesEditor!.state.doc.toString();
    if (!content || isSponsor) {
        const resp = await fetch(
            `/api/notes/${hole}/${getLang()}`,
            content ? { method: 'PUT', body: content} : { method: 'DELETE' },
        );
        if (resp.status !== 204) $<HTMLButtonElement>('#upsert-notes-btn').disabled = false;
        else holeLangNotesContent = content;
    }
};

function parseSubstitutions() {
    const { value } = $<HTMLInputElement>('#notes-substitutions');
    localStorage.setItem(`${getLang()}-substitutions`, value);
    const pattern = /s\/((?:[^\/]|\\\/)*)\/((?:[^\/]|\\\/)*)\/([dgimsuvy]*)/g;
    substitutions = [...value.matchAll(pattern)]
        .map(match => (
            {pattern: new RegExp(match[1], [...new Set(match[3])].join('')), replacement: JSON.parse(`"${match[2]}"`)}
        ));
    $<HTMLButtonElement>('#convert-notes-btn').disabled = substitutions.length < 1;
}

async function convertNotesAndRun(): Promise<boolean> {
    if (editor && holeLangNotesEditor) {
        let newCode = holeLangNotesEditor.state.doc.toString();
        for (const {pattern, replacement} of substitutions) {
            newCode = newCode.replace(pattern, replacement);
        }
        const code = editor.state.doc.toString();
        if (code !== newCode) {
            setState(newCode, editor);
            return await submit(editor, updateReadonlyPanels);
        }
    }
    return false;
}

layout.registerComponentFactoryFunction('holeLangNotes', async container => {
    container.setTitle(getTitle('holeLangNotes'));
    autoFocus(container);

    const header = (<header>
        <div>
            <input id="notes-substitutions" placeholder="Perl-like s/// expressions" size="50"></input>
            <button class="btn blue" id="convert-notes-btn" disabled>Convert & Run</button>
        </div>
        <button class="btn blue" id="upsert-notes-btn" disabled>Save</button>
    </header>) as HTMLElement;
    const editorDiv = <div id="holeLangNotesEditor"></div> as HTMLDivElement;

    container.element.id = 'holeLangNotesEditor-section';
    container.element.append(editorDiv, header);

    await afterDOM();
    makeHoleLangNotesEditor(editorDiv);

    $('#upsert-notes-btn').onclick = upsertNotes;
    $('#convert-notes-btn').onclick = convertNotesAndRun;
    $('#notes-substitutions').oninput = parseSubstitutions;
    $('#holeLangNotesEditor-section').onkeydown = ctrlEnter(async () => {
        if (await convertNotesAndRun()) {
            await upsertNotes();
        }
    });

    await afterDOM();
    updateHoleLangNotesContent();
    parseSubstitutions();
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
    container.setTitle(getTitle('scoreboard'));
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
    container.setTitle(getTitle('details'));
    autoFocus(container);
    const details = $<HTMLTemplateElement>('#template-details').content.cloneNode(true) as HTMLDetailsElement;
    container.element.append(details);
    container.element.id = 'details-content';
    initCopyButtons(container.element.querySelectorAll('[data-copy]'));
});

const titles: Record<string, string | undefined> = {
    details: 'Details',
    scoreboard: 'Scoreboard',
    exp: 'Expected',
    out: 'Output',
    err: 'Errors',
    arg: 'Arguments',
    diff: 'Diff',
    code: 'Code',
    langWiki: 'Wiki',
    holeLangNotes: 'Notes',
};

function getTitle(name: string) {
    return titles[name] ?? name;
}

function plainComponent(componentType: string): ComponentItemConfig {
    return {
        type: 'component',
        componentType,
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

function getPoolFromLayoutConfig(config: LayoutConfig) {
    const allComponents = ['code', 'scoreboard', 'arg', 'exp', 'out', 'err', 'diff', 'details', 'langWiki'];
    if (isSponsor || hasNotes) allComponents.push('holeLangNotes');
    const activeComponents = getComponentNamesRecursive(config.root!);
    return allComponents.filter(x => !activeComponents.includes(x));
}

function getComponentNamesRecursive(config: RootItemConfig): string[] {
    if (config.type === 'component'){
        return [config.componentType as string];
    }
    return config.content.flatMap(getComponentNamesRecursive);
}

const defaultViewState: ViewState = {
    version: 1,
    config: defaultLayout,
    isWide: false,
};

interface ViewState {
    version: 1;
    config: ResolvedLayoutConfig | LayoutConfig;
    isWide: boolean;
}

function getViewState(): ViewState {
    return {
        version: 1,
        config: layout.saveLayout(),
        isWide,
    };
}

const saveLayout = throttle(() => {
    const state = getViewState();
    if (!state.config.root) return;
    localStorage.setItem('lastViewState', JSON.stringify(state));
}, 2000);


async function applyInitialLayout() {
    const saved = localStorage.getItem('lastViewState');
    const viewState = saved
        ? JSON.parse(saved) as ViewState
        : defaultViewState;
    await applyViewState(viewState);
}

async function applyViewState(viewState: ViewState) {
    applyingDefault = true;
    toggleMobile(false);
    Object.keys(poolElements).map(removePoolItem);
    setWide(viewState.isWide);
    setTall(false);
    let { config } = viewState;
    if (LayoutConfig.isResolved(config))
        config = LayoutConfig.fromResolved(config);
    layout.loadLayout(config);
    getPoolFromLayoutConfig(config).forEach(addPoolItem);
    await afterDOM();
    checkMobile();
    applyingDefault = false;
}

applyInitialLayout();

/**
 * Try to add after selected item, with sensible defaults
 */
function addItemFromPool(componentName: string) {
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

$('#revert-layout').addEventListener('click', () => applyViewState(defaultViewState));

function setWide(b: boolean) {
    isWide = b;
    document.documentElement.classList.toggle('full-width', b);
}

function setTall(b: boolean) {
    document.documentElement.classList.toggle('full-height', b);
}

$('#make-wide').addEventListener('click', () => setWide(true));
$('#make-narrow').addEventListener('click', () => setWide(false));

$('#make-tall').addEventListener('click', () => setTall(true));
$('#make-short').addEventListener('click', () => setTall(false));

function addPoolItem(componentType: string) {
    poolElements[componentType]?.remove();
    const el = (<span class="btn">{getTitle(componentType)}</span>);
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
        addPoolItem(target.componentType as string);
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
    saveLayout();
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
        item.type = 'column';
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
    const currLayout = layout.saveLayout();
    if (currLayout.root) {
        mutateDeep(currLayout.root as DeepMutable<ResolvedRootItemConfig>, isMobile);
        layout.loadLayout(LayoutConfig.fromResolved(currLayout));
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
