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

function pushInterpolationStack(state) {
    (state.interpolationStack || (state.interpolationStack = [])).push(state.tokenize);
}

function popInterpolationStack(state) {
    return (state.interpolationStack || (state.interpolationStack = [])).pop();
}

function sizeInterpolationStack(state) {
    return state.interpolationStack ? state.interpolationStack.length : 0;
}

function tokenDartString(quote, stream, state, raw) {
    let tripleQuoted = false;
    if (stream.eat(quote)) {
        if (stream.eat(quote)) tripleQuoted = true;
        else return 'string'; //empty string
    }
    function tokenStringHelper(stream, state) {
        let escaped = false;
        while (!stream.eol()) {
            if (!raw && !escaped && stream.peek() == '$') {
                pushInterpolationStack(state);
                state.tokenize = tokenInterpolation;
                return 'string';
            }
            const next = stream.next();
            if (next == quote && !escaped && (!tripleQuoted || stream.match(quote + quote))) {
                state.tokenize = null;
                break;
            }
            escaped = !raw && !escaped && next == '\\';
        }
        return 'string';
    }
    state.tokenize = tokenStringHelper;
    return tokenStringHelper(stream, state);
}

function tokenInterpolation(stream, state) {
    stream.eat('$');
    if (stream.eat('{')) {
    // let clike handle the content of ${...},
    // we take over again when "}" appears (see hooks).
        state.tokenize = null;
    }
    else {
        state.tokenize = tokenInterpolationIdentifier;
    }
    return null;
}

function tokenInterpolationIdentifier(stream, state) {
    stream.eatWhile(/[\w_]/);
    state.tokenize = popInterpolationStack(state);
    return 'variable';
}

export const dart = clike({
    keywords: words('this super static final const abstract class extends external factory ' +
                  'implements mixin get native set typedef with enum throw rethrow ' +
                  'assert break case continue default in return new deferred async await covariant ' +
                  'try catch finally do else for if switch while import library export ' +
                  'part of show hide is as extension on yield late required'),
    blockKeywords: words('try catch finally do else for if switch while'),
    builtin: words('void bool num int double dynamic var String Null Never'),
    atoms: words('true false null'),
    hooks: {
        '@': function (stream) {
            stream.eatWhile(/[\w\$_\.]/);
            return 'meta';
        },

        // custom string handling to deal with triple-quoted strings and string interpolation
        "'": function (stream, state) {
            return tokenDartString("'", stream, state, false);
        },
        '"': function (stream, state) {
            return tokenDartString('"', stream, state, false);
        },
        'r': function (stream, state) {
            const peek = stream.peek();
            if (peek == "'" || peek == '"') {
                return tokenDartString(stream.next(), stream, state, true);
            }
            return false;
        },

        '}': function (_stream, state) {
            // "}" is end of interpolation, if interpolation stack is non-empty
            if (sizeInterpolationStack(state) > 0) {
                state.tokenize = popInterpolationStack(state);
                return null;
            }
            return false;
        },

        '/': function (stream, state) {
            if (!stream.eat('*')) return false;
            state.tokenize = tokenNestedComment(1);
            return state.tokenize(stream, state);
        },
        'token': function (stream, _, style) {
            if (style == 'variable') {
                // Assume uppercase symbols are classes
                const isUpper = RegExp('^[_$]*[A-Z][a-zA-Z0-9_$]*$','g');
                if (isUpper.test(stream.current())) {
                    return 'type';
                }
            }
        },
    },
});
