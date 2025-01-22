import { StyleModule } from 'style-mod'; // Already used by CodeMirror

/* eslint no-unused-vars: ["off"] */

const shadowStyle = new StyleModule({
    ':host': {
        display: 'inline-block',
        position: 'relative',
        verticalAlign: 'middle',
        width: '1ch',
        height: '1em',
        // Omit emoji because Apple Color Emoji has a full-width glyph for non-emoji '0'. (???)
        fontFamily: "'Source Code Pro', monospace",
        lineHeight: 0.8,
        cursor: 'text',
    },
    ':host > span': {
        'WebkitTextFillColor': 'transparent',
        '&:before, &:after': {
            WebkitTextFillColor: 'currentcolor',
            fontSize: '70%',
            position: 'absolute',
            pointerEvents: 'none',
        },
        '&:before': {
            content: 'attr(data-h)',
            left: 0,
            top: 0,
        },
        '&:after': {
            content: 'attr(data-l)',
            right: 0,
            bottom: 0,
        },
    },

    ':host([c])': {
        pointerEvents: 'none',
    },
    ':host([c]) > span': {
        pointerEvents: 'auto',
    },
});

// <u-p>&#...;</u-p> renders a single character and allows selection and copying.
// <u-p c="&#...;"></u-p> renders a single character but doesn't allow copying.
//
// TODO:
// Title doesn't show up for the attribute use case in Safari and others.
// Firefox doesn't (yet) let the selection cross shadow root boundaries. (cf. bug #1867058)
export default class UnprintableElement extends HTMLElement {
    #span;

    constructor(text = '') {
        super();

        const shadow = this.attachShadow({ mode: 'closed' });
        StyleModule.mount(shadow, shadowStyle);

        this.#span = document.createElement('span');
        shadow.append(this.#span);

        if (text) this.textContent = text;
    }

    // For the simplicity, we only update the shadow DOM when connected.
    connectedCallback() {
        let c = this.getAttribute('c'), h, l, t;
        const ignoreTextContent = !!c;
        c = c || this.textContent || '';
        if (c.length == 0) {
            h = '⌜';
            l = '⌟';
            t = '(empty)';
        }
        else if (c.length == 1) {
            const code = c.charCodeAt(0);
            h = '0123456789ABCDEF'[code / 16 | 0];
            l = '0123456789ABCDEF'[code % 16];
            t = '\\u' + code.toString(16);
        }
        else {
            h = '+';
            l = '+';
            c = '';
            t = '(multiple)';
        }
        this.#span.setAttribute('data-h', h);
        this.#span.setAttribute('data-l', l);
        this.#span.textContent = ignoreTextContent ? '' : c;
        this.#span.title = t;
    }

    static PATTERN = /([\x00-\x08\x0B-\x1F\x7F-\xA0])/g;

    static escape(text: string): DocumentFragment {
        const frag = document.createDocumentFragment();
        const parts = text.split(UnprintableElement.PATTERN);
        for (let i = 0; i < parts.length - 1; i += 2) {
            frag.append(parts[i], new UnprintableElement(parts[i + 1]));
        }
        frag.append(parts[parts.length - 1]);
        return frag;
    }
}

customElements.define('u-p', UnprintableElement);

declare global {
    namespace JSX {
        interface IntrinsicElements {
            ['u-p']: any;
        }
    }
}
