import tsParser from '@typescript-eslint/parser';
import tsPlugin from '@typescript-eslint/eslint-plugin';

// Debug which files are matched with:
// DEBUG=eslint:linter node_modules/.bin/eslint js 2>&1 | grep 'Linting code'

export default [
    { ignores: ['js/vendor/**'] }, // Globally ignore vendored code.
    {
        files:           ['js/**/*.ts', 'js/**/*.tsx'],
        languageOptions: { parser: tsParser },
        plugins:         { '@typescript-eslint': tsPlugin },
        rules: {
            'array-bracket-newline':       ['error', 'consistent'],
            'arrow-parens':                ['error', 'as-needed'],
            'brace-style':                 ['error', 'stroustrup'],
            'camelcase':                   ['error'],
            'comma-dangle':                ['error', 'always-multiline'],
            'indent':                      ['error'],
            'keyword-spacing':             ['error'],
            'no-duplicate-imports':        ['error'],
            'no-trailing-spaces':          ['error'],
            'no-unused-vars':              ['error'],
            'no-useless-assignment':       ['error'],
            'no-var':                      ['error'],
            'prefer-const':                ['error'],
            'prefer-destructuring':        ['error', { object: true }],
            'quote-props':                 ['error', 'consistent-as-needed'],
            'quotes':                      ['error', 'single', { avoidEscape: true }],
            'semi':                        ['error', 'always'],
            'space-before-function-paren': ['error', { named: 'never' }],
        },
    },
];
