// From https://github.com/codemirror/legacy-modes/blob/main/mode/clike.js
// Modified to support multiline strings in C.
function Context(indented, column, type, info, align, prev) {
    this.indented = indented;
    this.column = column;
    this.type = type;
    this.info = info;
    this.align = align;
    this.prev = prev;
}
function pushContext(state, col, type, info) {
    let indent = state.indented;
    if (state.context && state.context.type == 'statement' && type != 'statement')
        indent = state.context.indented;
    return state.context = new Context(indent, col, type, info, null, state.context);
}
function popContext(state) {
    const t = state.context.type;
    if (t == ')' || t == ']' || t == '}')
        state.indented = state.context.indented;
    return state.context = state.context.prev;
}

function typeBefore(stream, state, pos) {
    if (state.prevToken == 'variable' || state.prevToken == 'type') return true;
    if (/\S(?:[^- ]>|[*\]])\s*$|\*$/.test(stream.string.slice(0, pos))) return true;
    if (state.typeAtEndOfLine && stream.column() == stream.indentation()) return true;
}

function isTopScope(context) {
    for (;;) {
        if (!context || context.type == 'top') return true;
        if (context.type == '}' && context.prev.info != 'namespace') return false;
        context = context.prev;
    }
}

export function clike(parserConfig) {
    const statementIndentUnit = parserConfig.statementIndentUnit,
        dontAlignCalls = parserConfig.dontAlignCalls,
        keywords = parserConfig.keywords || {},
        types = parserConfig.types || {},
        builtin = parserConfig.builtin || {},
        blockKeywords = parserConfig.blockKeywords || {},
        defKeywords = parserConfig.defKeywords || {},
        atoms = parserConfig.atoms || {},
        hooks = parserConfig.hooks || {},
        multiLineStrings = parserConfig.multiLineStrings,
        indentStatements = parserConfig.indentStatements !== false,
        indentSwitch = parserConfig.indentSwitch !== false,
        namespaceSeparator = parserConfig.namespaceSeparator,
        isPunctuationChar = parserConfig.isPunctuationChar || /[\[\]{}\(\),;\:\.]/,
        numberStart = parserConfig.numberStart || /[\d\.]/,
        number = parserConfig.number || /^(?:0x[a-f\d]+|0b[01]+|(?:\d+\.?\d*|\.\d+)(?:e[-+]?\d+)?)(u|ll?|l|f)?/i,
        isOperatorChar = parserConfig.isOperatorChar || /[+\-*&%=<>!?|\/]/,
        isIdentifierChar = parserConfig.isIdentifierChar || /[\w\$_\xa1-\uffff]/,
        // An optional function that takes a {string} token and returns true if it
        // should be treated as a builtin.
        isReservedIdentifier = parserConfig.isReservedIdentifier || false;

    let curPunc, isDefKeyword;

    function tokenBase(stream, state) {
        const ch = stream.next();
        if (hooks[ch]) {
            const result = hooks[ch](stream, state);
            if (result !== false) return result;
        }
        if (ch == '"' || ch == "'") {
            state.tokenize = tokenString(ch);
            return state.tokenize(stream, state);
        }
        if (numberStart.test(ch)) {
            stream.backUp(1);
            if (stream.match(number)) return 'number';
            stream.next();
        }
        if (isPunctuationChar.test(ch)) {
            curPunc = ch;
            return null;
        }
        if (ch == '/') {
            if (stream.eat('*')) {
                state.tokenize = tokenComment;
                return tokenComment(stream, state);
            }
            if (stream.eat('/')) {
                stream.skipToEnd();
                return 'comment';
            }
        }
        if (isOperatorChar.test(ch)) {
            while (!stream.match(/^\/[\/*]/, false) && stream.eat(isOperatorChar)) {}
            return 'operator';
        }
        stream.eatWhile(isIdentifierChar);
        if (namespaceSeparator) while (stream.match(namespaceSeparator))
            stream.eatWhile(isIdentifierChar);

        const cur = stream.current();
        if (contains(keywords, cur)) {
            if (contains(blockKeywords, cur)) curPunc = 'newstatement';
            if (contains(defKeywords, cur)) isDefKeyword = true;
            return 'keyword';
        }
        if (contains(types, cur)) return 'type';
        if (contains(builtin, cur)
        || (isReservedIdentifier && isReservedIdentifier(cur))) {
            if (contains(blockKeywords, cur)) curPunc = 'newstatement';
            return 'builtin';
        }
        if (contains(atoms, cur)) return 'atom';
        return 'variable';
    }

    function tokenString(quote) {
        return function (stream, state) {
            let escaped = false, next, end = false;
            while ((next = stream.next()) != null) {
                if (next == quote && !escaped) {end = true; break}
                escaped = !escaped && next == '\\';
            }
            if (end || !(escaped || multiLineStrings))
                state.tokenize = null;
            return 'string';
        };
    }

    function tokenComment(stream, state) {
        let maybeEnd = false, ch;
        while (ch = stream.next()) {
            if (ch == '/' && maybeEnd) {
                state.tokenize = null;
                break;
            }
            maybeEnd = (ch == '*');
        }
        return 'comment';
    }

    function maybeEOL(stream, state) {
        if (parserConfig.typeFirstDefinitions && stream.eol() && isTopScope(state.context))
            state.typeAtEndOfLine = typeBefore(stream, state, stream.pos);
    }

    // Interface

    return {
        startState: function (indentUnit) {
            return {
                tokenize: null,
                context: new Context(-indentUnit, 0, 'top', null, false),
                indented: 0,
                startOfLine: true,
                prevToken: null,
            };
        },

        token: function (stream, state) {
            let ctx = state.context;
            if (stream.sol()) {
                if (ctx.align == null) ctx.align = false;
                state.indented = stream.indentation();
                state.startOfLine = true;
            }
            if (stream.eatSpace()) { maybeEOL(stream, state); return null }
            curPunc = isDefKeyword = null;
            let style = (state.tokenize || tokenBase)(stream, state);
            if (style == 'comment' || style == 'meta') return style;
            if (ctx.align == null) ctx.align = true;

            if (curPunc == ';' || curPunc == ':' || (curPunc == ',' && stream.match(/^\s*(?:\/\/.*)?$/, false)))
                while (state.context.type == 'statement') popContext(state);
            else if (curPunc == '{') pushContext(state, stream.column(), '}');
            else if (curPunc == '[') pushContext(state, stream.column(), ']');
            else if (curPunc == '(') pushContext(state, stream.column(), ')');
            else if (curPunc == '}') {
                while (ctx.type == 'statement') ctx = popContext(state);
                if (ctx.type == '}') ctx = popContext(state);
                while (ctx.type == 'statement') ctx = popContext(state);
            }
            else if (curPunc == ctx.type) popContext(state);
            else if (indentStatements &&
               (((ctx.type == '}' || ctx.type == 'top') && curPunc != ';') ||
                (ctx.type == 'statement' && curPunc == 'newstatement'))) {
                pushContext(state, stream.column(), 'statement', stream.current());
            }

            if (style == 'variable' &&
          ((state.prevToken == 'def' ||
            (parserConfig.typeFirstDefinitions && typeBefore(stream, state, stream.start) &&
             isTopScope(state.context) && stream.match(/^\s*\(/, false)))))
                style = 'def';

            if (hooks.token) {
                const result = hooks.token(stream, state, style);
                if (result !== undefined) style = result;
            }

            if (style == 'def' && parserConfig.styleDefs === false) style = 'variable';

            state.startOfLine = false;
            state.prevToken = isDefKeyword ? 'def' : style || curPunc;
            maybeEOL(stream, state);
            return style;
        },

        indent: function (state, textAfter, context) {
            if (state.tokenize != tokenBase && state.tokenize != null || state.typeAtEndOfLine) return null;
            let ctx = state.context;
            const firstChar = textAfter && textAfter.charAt(0);
            const closing = firstChar == ctx.type;
            if (ctx.type == 'statement' && firstChar == '}') ctx = ctx.prev;
            if (parserConfig.dontIndentStatements)
                while (ctx.type == 'statement' && parserConfig.dontIndentStatements.test(ctx.info))
                    ctx = ctx.prev;
            if (hooks.indent) {
                const hook = hooks.indent(state, ctx, textAfter, context.unit);
                if (typeof hook == 'number') return hook;
            }
            const switchBlock = ctx.prev && ctx.prev.info == 'switch';
            if (parserConfig.allmanIndentation && /[{(]/.test(firstChar)) {
                while (ctx.type != 'top' && ctx.type != '}') ctx = ctx.prev;
                return ctx.indented;
            }
            if (ctx.type == 'statement')
                return ctx.indented + (firstChar == '{' ? 0 : statementIndentUnit || context.unit);
            if (ctx.align && (!dontAlignCalls || ctx.type != ')'))
                return ctx.column + (closing ? 0 : 1);
            if (ctx.type == ')' && !closing)
                return ctx.indented + (statementIndentUnit || context.unit);

            return ctx.indented + (closing ? 0 : context.unit) +
        (!closing && switchBlock && !/^(?:case|default)\b/.test(textAfter) ? context.unit : 0);
        },

        languageData: {
            indentOnInput: indentSwitch ? /^\s*(?:case .*?:|default:|\{\}?|\})$/ : /^\s*[{}]$/,
            commentTokens: {line: '//', block: {open: '/*', close: '*/'}},
            autocomplete: Object.keys(keywords).concat(Object.keys(types)).concat(Object.keys(builtin)).concat(Object.keys(atoms)),
            ...parserConfig.languageData,
        },
    };
};

function words(str) {
    const obj = {}, words = str.split(' ');
    for (let i = 0; i < words.length; ++i) obj[words[i]] = true;
    return obj;
}
function contains(words, word) {
    if (typeof words === 'function') {
        return words(word);
    }
    else {
        return words.propertyIsEnumerable(word);
    }
}
const cKeywords = 'auto if break case register continue return default do sizeof ' +
    'static else struct switch extern typedef union for goto while enum const ' +
    'volatile inline restrict asm fortran';

// Do not use this. Use the cTypes function below. This is global just to avoid
// excessive calls when cTypes is being called multiple times during a parse.
const basicCTypes = words('int long char short double float unsigned signed ' +
                        'void bool');

// Returns true if identifier is a "C" type.
// C type is defined as those that are reserved by the compiler (basicTypes),
// and those that end in _t (Reserved by POSIX for types)
// http://www.gnu.org/software/libc/manual/html_node/Reserved-Names.html
function cTypes(identifier) {
    return contains(basicCTypes, identifier) || /.+_t$/.test(identifier);
}

const cBlockKeywords = 'case do else for if switch while struct enum union';
const cDefKeywords = 'struct enum union';

function cppHook(stream, state) {
    if (!state.startOfLine) return false;
    let next = null;
    for (let ch; ch = stream.peek();) {
        if (ch == '\\' && stream.match(/^.$/)) {
            next = cppHook;
            break;
        }
        else if (ch == '/' && stream.match(/^\/[\/\*]/, false)) {
            break;
        }
        stream.next();
    }
    state.tokenize = next;
    return 'meta';
}

function pointerHook(_stream, state) {
    if (state.prevToken == 'type') return 'type';
    return false;
}

// For C and C++ (and ObjC): identifiers starting with __
// or _ followed by a capital letter are reserved for the compiler.
function cIsReservedIdentifier(token) {
    if (!token || token.length < 2) return false;
    if (token[0] != '_') return false;
    return (token[1] == '_') || (token[1] !== token[1].toLowerCase());
}

// C#-style strings where "" escapes a quote.
function tokenAtString(stream, state) {
    let next;
    while ((next = stream.next()) != null) {
        if (next == '"' && !stream.eat('"')) {
            state.tokenize = null;
            break;
        }
    }
    return 'string';
}

export const c = clike({
    keywords: words(cKeywords),
    types: cTypes,
    blockKeywords: words(cBlockKeywords),
    defKeywords: words(cDefKeywords),
    typeFirstDefinitions: true,
    atoms: words('NULL true false'),
    isReservedIdentifier: cIsReservedIdentifier,
    hooks: {
        '#': cppHook,
        '*': pointerHook,
    },
    multiLineStrings: true,
});

export const csharp = clike({
    keywords: words('abstract as async await base break case catch checked class const continue' +
                  ' default delegate do else enum event explicit extern finally fixed for' +
                  ' foreach goto if implicit in interface internal is lock namespace new' +
                  ' operator out override params private protected public readonly ref return sealed' +
                  ' sizeof stackalloc static struct switch this throw try typeof unchecked' +
                  ' unsafe using virtual void volatile while add alias ascending descending dynamic from get' +
                  ' global group into join let orderby partial remove select set value var yield'),
    types: words('Action Boolean Byte Char DateTime DateTimeOffset Decimal Double Func' +
               ' Guid Int16 Int32 Int64 Object SByte Single String Task TimeSpan UInt16 UInt32' +
               ' UInt64 bool byte char decimal double short int long object'  +
               ' sbyte float string ushort uint ulong'),
    blockKeywords: words('catch class do else finally for foreach if struct switch try while'),
    defKeywords: words('class interface namespace struct var'),
    typeFirstDefinitions: true,
    atoms: words('true false null'),
    hooks: {
        '@': function (stream, state) {
            if (stream.eat('"')) {
                state.tokenize = tokenAtString;
                return tokenAtString(stream, state);
            }
            stream.eatWhile(/[\w\$_]/);
            return 'meta';
        },
    },
});
