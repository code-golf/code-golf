export default {
    js2svg: { finalNewline: true  },
    plugins: [
        {
            name: 'convertTransform',
            active: false,
        },
        {
            name: 'preset-default',
            // We need top-level IDs for <use href="#foo"/> to work.
            params: { overrides: { cleanupIds: false } },
        },
    ],
};
