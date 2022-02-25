
export const comma = i => i.toLocaleString('en');

// Adapted from https://codegolf.stackexchange.com/a/119563.
export const ord = i => [, 'st', 'nd', 'rd'][i % 100 >> 3 ^ 1 && i % 10] || 'th';

// charLen adapted from https://mths.be/punycode.
export const byteLen = str => new TextEncoder().encode(str).length;
export const charLen = str => {
    let i = 0, len = 0;

    while (i < str.length) {
        const value = str.charCodeAt(i++);

        if (value >= 0xD800 && value <= 0xDBFF && i < str.length) {
            // It's a high surrogate, and there is a next character.
            const extra = str.charCodeAt(i++);

            // Low surrogate.
            if ((extra & 0xFC00) == 0xDC00) {
                len++;
            }
            else {
                // It's an unmatched surrogate; only append this code unit, in
                // case the next code unit is the high surrogate of a
                // surrogate pair.
                len++;
                i--;
            }
        }
        else {
            len++;
        }
    }

    return len;
};
