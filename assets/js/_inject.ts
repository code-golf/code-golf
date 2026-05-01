// Exports injected into every TS/TSX file by default.

// Our default JSX factory.
export function createElement(
    tag: string,
    attrs: {[key: string]: string},
    ...children: (Node | string)[]
) {
    const element = tag === 'svg' || tag === 'use'
        ? document.createElementNS('http://www.w3.org/2000/svg', tag)
        : document.createElement(tag);

    // Set all defined/non-null attributes.
    Object.entries(attrs ?? {})
        .filter(([     , value]) => value != null)
        .forEach(([name, value]) => element.setAttribute(name, value));

    // Set all defined/non-null children.
    element.append(...children.flat().filter(e => e != null));

    return element;
}
