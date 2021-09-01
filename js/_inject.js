// Exports injected into every JS file by default.

// Our default JSX factory.
export function createElement(tag, attrs, ...children) {
    const element = document.createElement(tag);

    Object.entries(attrs ?? {}).forEach(
        ([name, value]) => element.setAttribute(name, value));

    element.append(...children.flat());

    return element;
}
