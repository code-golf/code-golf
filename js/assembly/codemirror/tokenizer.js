import {
    fetchMnemonic, isSizeHint, prefixes, isRegister,
    isDirective, scanIdentifier
} from '@defasm/core';
import { ContextTracker, ExternalTokenizer, InputStream } from '@lezer/lr';

import * as Terms from './parser.terms.js';

var tok, pureString;

/** @param {InputStream} input */
function next(input)
{
    tok = '';
    let char;
    
    while(input.next >= 0 && input.next != 10 && String.fromCharCode(input.next).match(/\s/))
        input.advance();
    if(input.next >= 0 && !(char = String.fromCharCode(input.next)).match(/[.$\w]/))
    {
        tok = char;
        input.advance();
        pureString = false;
    }
    else
        while(input.next >= 0 && (char = String.fromCharCode(input.next)).match(/[.$\w]/))
        {
            tok += char;
            input.advance();
        }
        
    tok = tok.toLowerCase() || '\n';

    return tok;
}

/** @param {InputStream} input */
function peekNext(input)
{
    let i = 0, char;
    while((char = input.peek(i)) >= 0 && char != 10 && String.fromCharCode(char).match(/\s/))
        i++;
    if((char = input.peek(i)) >= 0 && !(char = String.fromCharCode(char)).match(/[.$\w]/))
        return char;
    
    let result = '';
    while((char = input.peek(i)) >= 0 && (char = String.fromCharCode(char)).match(/[.$\w]/))
    {
        result += char;
        i++;
    }
    return result.toLowerCase() || '\n';
}

const
    STATE_SYNTAX_INTEL = 1,
    STATE_SYNTAX_PREFIX = 2,
    STATE_IN_INSTRUCTION = 4,
    STATE_ALLOW_IMM = 8,
    STATE_SYNTAX_X86 = 16;

/** @param {import('@defasm/core/parser.js').Syntax} initialSyntax */
export const ctxTracker = (initialSyntax, bitness) => new ContextTracker({
    start:
        (initialSyntax.intel * STATE_SYNTAX_INTEL) |
        (initialSyntax.prefix * STATE_SYNTAX_PREFIX) |
        ((bitness == 32) * STATE_SYNTAX_X86),
    shift: (ctx, term, stack, input) => {
        if(term == Terms.Opcode)
            ctx |= STATE_IN_INSTRUCTION | STATE_ALLOW_IMM;
        else if(term == Terms.RelOpcode || term == Terms.IOpcode
            || term == Terms.IRelOpcode || term == Terms.symEquals
            || term == Terms.Directive)
            ctx |= STATE_IN_INSTRUCTION;
        else if((ctx & STATE_IN_INSTRUCTION) && term != Terms.Space)
        {
            if(input.next == ','.charCodeAt(0))
                ctx |= STATE_ALLOW_IMM;
            else
                ctx &= ~STATE_ALLOW_IMM;
        }
        
        if(input.next == '\n'.charCodeAt(0) || input.next == ';'.charCodeAt(0))
            ctx &= ~STATE_IN_INSTRUCTION;
        if(term != Terms.Directive)
            return ctx;
        let result = ctx, syntax = next(input);
        if(syntax == ".intel_syntax")
        {
            result |= STATE_SYNTAX_INTEL;
            result &= ~STATE_SYNTAX_PREFIX;
        }
        else if(syntax == ".att_syntax")
        {
            result &= ~STATE_SYNTAX_INTEL;
            result |= STATE_SYNTAX_PREFIX;
        }
        else
            return ctx;
        const pref = next(input);
        if(pref == 'prefix')
            result |= STATE_SYNTAX_PREFIX;
        else if(pref == 'noprefix')
            result &= ~STATE_SYNTAX_PREFIX;
        else if(pref != '\n' && pref != ';' && ((result & STATE_SYNTAX_INTEL) || pref != '#'))
            return ctx;
        
        return result;
    },
    hash: ctx => ctx,
    strict: false
});

/** @param {InputStream} input */
function tokenize(ctx, input)
{
    const intel = ctx & STATE_SYNTAX_INTEL,
          prefix = ctx & STATE_SYNTAX_PREFIX,
          bitness = (ctx & STATE_SYNTAX_X86) ? 32 : 64;
    if(tok == (intel ? ';' : '#'))
    {
        while(input.next >= 0 && input.next != '\n'.charCodeAt(0))
            input.advance();
        return Terms.Comment;
    }

    if(!(ctx & STATE_IN_INSTRUCTION))
    {
        const nextTok = peekNext(input);
        if(nextTok == '=' || nextTok == ':' || intel && (nextTok == 'equ' || isDirective(nextTok, true)))
            return Terms.SymbolName;

        if(tok == '%' && intel)
            return isDirective('%' + next(input), true) ? Terms.Directive : null;
        
        if(isDirective(tok, intel))
            return Terms.Directive;

        if(intel && tok == 'offset')
            return Terms.Offset;

        if(prefixes.hasOwnProperty(tok))
            return Terms.Prefix;
        
        if(tok == '=' || intel && tok == 'equ')
            return Terms.symEquals;

        let opcode = tok, interps = fetchMnemonic(opcode, intel, !intel, bitness);
        if(interps.length > 0)
            return interps[0].relative
            ?
                intel ? Terms.IRelOpcode : Terms.RelOpcode
            :
                intel ? Terms.IOpcode : Terms.Opcode;
        return null;
    }
    if((ctx & STATE_ALLOW_IMM) && tok[0] == '$')
    {
        input.pos -= tok.length - 1;
        return Terms.immPrefix;
    }

    if(tok == '@')
    {
        next(input);
        return Terms.SpecialWord;
    }

    if(tok == '%' && prefix)
        return isRegister(next(input), bitness) ? Terms.Register : null;
    
    if(tok == '{')
    {
        if((!prefix || next(input) == '%') && isRegister(next(input), bitness))
            return null;
        while(tok != '\n' && tok != '}')
            next(input);
        return Terms.VEXRound;
    }
    
    if(intel && isSizeHint(tok))
    {
        let prevEnd = input.pos;
        if(",;\n{:".includes(next(input)))
        {
            input.pos = prevEnd;
            return Terms.word;
        }

        if(tok == 'ptr')
        {
            let nextPrevEnd = input.pos;
            input.pos = ",;\n{:".includes(next(input)) ? prevEnd : nextPrevEnd;
            return Terms.Ptr;
        }

        input.pos = prevEnd;
        return Terms.Ptr;
    }

    const idType = scanIdentifier(tok, intel);
    if(idType === null)
        return null;
    if(!prefix && isRegister(tok, bitness))
        return Terms.Register;
    if(idType == 'symbol')
        return Terms.word;
    return Terms.number;
}
export const tokenizer = new ExternalTokenizer(
    (input, stack) => {
        if(input.next < 0 || String.fromCharCode(input.next).match(/\s/))
            return;

        pureString = true;
        next(input);
        const type = tokenize(stack.context, input);
        if(type !== null || pureString)
            input.acceptToken(type ?? Terms.None);
        
    }, {
        contextual: false
    }
);