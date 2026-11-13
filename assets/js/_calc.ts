// A Python 3 calculator: integer arithmetic is arbitrary precision (bigint),
// "/" gives a float. Never evals; it's a small precedence climbing parser.
// "%%" is an extra: the truncating remainder of C/C#/Java/JS.

type Val = bigint | number;

// Upper bound on the bit size of "**" and "<<" results (~1,200 digits).
const MAX_BITS = 4096;

const fail = (msg: string): never => { throw new Error(msg); };

const int = (v: Val) => typeof v == 'bigint' ? v : fail('Integer required');

const nonZero = (v: Val) => v == 0 ? fail('Division by zero') : v;

const shiftCount = (v: Val) => int(v) < 0n ? fail('Negative shift count') : int(v);

// Apply an integer op to two ints, else a float op to their float values.
const num = (i: (a: bigint, b: bigint) => Val, f: (a: number, b: number) => number) =>
    (a: Val, b: Val) => typeof a == 'bigint' && typeof b == 'bigint' ? i(a, b) : f(Number(a), Number(b));

const floorDiv = (a: bigint, b: bigint) => {
    const q = a / b;
    return q * b != a && (a < 0n) != (b < 0n) ? q - 1n : q;
};

const floorMod = (a: number, b: number) => {
    const r = a % b;
    return r && (r < 0) != (b < 0) ? r + b : r;
};

const pow = (a: Val, b: Val): Val => {
    if (typeof a == 'bigint' && typeof b == 'bigint' && b >= 0n) {
        const abs = a < 0n ? -a : a;
        if (abs > 1n && b > 1n && Number(b) * Math.log2(Number(abs)) > MAX_BITS)
            fail('Result too large');
        return a ** b;
    }
    if (a == 0 && b < 0) fail('Division by zero');
    return Number(a) ** Number(b);
};

// Python's precedence, loosest first. Unary ops and "**" bind tighter.
const binary: { [op: string]: [number, (a: Val, b: Val) => Val] } = {
    '|':  [1, (a, b) => int(a) | int(b)],
    '^':  [2, (a, b) => int(a) ^ int(b)],
    '&':  [3, (a, b) => int(a) & int(b)],
    '<<': [4, (a, b) => shiftCount(b) > MAX_BITS ? fail('Result too large') : int(a) << int(b)],
    '>>': [4, (a, b) => int(a) >> shiftCount(b)],
    '+':  [5, num((a, b) => a + b, (a, b) => a + b)],
    '-':  [5, num((a, b) => a - b, (a, b) => a - b)],
    '*':  [6, num((a, b) => a * b, (a, b) => a * b)],
    '/':  [6, (a, b) => Number(a) / Number(nonZero(b))],
    '//': [6, (a, b) => num(floorDiv, (a, b) => Math.floor(a / b))(a, nonZero(b))],
    '%':  [6, (a, b) => num((a, b) => a - b * floorDiv(a, b), floorMod)(a, nonZero(b))],
    '%%': [6, (a, b) => num((a, b) => a % b, (a, b) => a % b)(a, nonZero(b))],
};

const unary: { [op: string]: (a: Val) => Val } = {
    '+': a => a,
    '-': a => -a,
    '~': a => ~int(a),
};

// Like Python's repr(): "4.0", "0.0001", "1e-05", "1e+16".
const formatFloat = (x: number) => {
    if (!isFinite(x)) fail('Result out of range');

    const [mantissa, exp] = x.toExponential().split('e');
    if (+exp < -4 || +exp >= 16)
        return `${mantissa}e${exp[0]}${exp.slice(1).padStart(2, '0')}`;

    return Number.isInteger(x) ? x.toFixed(1) : String(x);
};

export const calc = (src: string) => {
    const tokens = src.match(/\*\*|\/\/|%%|<<|>>|0[box][\da-f]+|(\d+\.?\d*|\.\d+)(e[+-]?\d+)?|\S/gi) ?? [];
    let i = 0;

    const atom = (): Val => {
        const token = tokens[i++];
        if (token == '(') {
            const value = expr();
            if (tokens[i++] != ')') fail('Expected )');
            return value;
        }
        if (/^\.?\d/.test(token))
            return /^0[box]|^[^.e]*$/i.test(token) ? BigInt(token) : Number(token);
        return fail(`Unexpected ${token ?? 'end'}`);
    };

    // "**" is right associative and its right side may be unary: 2**-1.
    const power = (): Val => {
        const base = atom();
        return tokens[i] == '**' ? (i++, pow(base, factor())) : base;
    };

    const factor = (): Val => {
        const op = unary[tokens[i]];
        return op ? (i++, op(factor())) : power();
    };

    const expr = (minPrec = 1): Val => {
        let lhs = factor();
        for (let op; (op = binary[tokens[i]]) && op[0] >= minPrec;) {
            i++;
            lhs = op[1](lhs, expr(op[0] + 1));
        }
        return lhs;
    };

    const value = expr();
    if (i < tokens.length) fail(`Unexpected ${tokens[i]}`);

    return typeof value == 'bigint' ? String(value) : formatFloat(value);
};
