// Exports injected into every JS file by default.

// Our default JSX factory.
export function createElement(tag, attrs, ...children) {
    const element = document.createElement(tag);

    // Set all defined attributes.
    Object.entries(attrs ?? {})
        .filter(([     , value]) => value !== null)
        .forEach(([name, value]) => element.setAttribute(name, value));

    // Set all defined children.
    element.append(...children.flat().filter(e => e !== null));

    return element;
}
