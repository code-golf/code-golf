module.exports = {
    js2svg: { finalNewline: true  },
    plugins: [
        {
            name: 'preset-default',
            // We need top-level IDs for <use href="#foo"/> to work.
            params: { overrides: { cleanupIds: false } },
        },
    ],
};
