// Exports injected into every JS file by default

// Small util functions.
/** Assume $ always succeeds and returns an HTMLElement */
export const $     = (selector: string) => document.querySelector(selector) as HTMLElement;
/** Assume $$ returns HTMLElements only */
export const $$    = (selector: string) => document.querySelectorAll(selector) as NodeListOf<HTMLElement>;
export const comma = (i: number) => i.toLocaleString('en');

// Our default JSX factory.
export function createElement(
    tag: string,
    attrs: {[key: string]: string},
    ...children: (Node | string)[]
) {
    const element = document.createElement(tag);

    // Set all defined/non-null attributes.
    Object.entries(attrs ?? {})
        .filter(([     , value]) => value != null)
        .forEach(([name, value]) => element.setAttribute(name, value));

    // Set all defined/non-null children.
    element.append(...children.flat().filter(e => e != null));

    return element;
}
