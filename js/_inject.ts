// Exports injected into every JS file by default

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
