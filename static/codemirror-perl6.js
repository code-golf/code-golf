// CodeMirror2 mode/perl/perl6.js (text/x-perl6)
// This is a part of CodeMirror from https://github.com/azawawi/farabi6 (ahmad.zawawi@gmail.com)
// This work is based on the perl.js mode by (TODO give credit)
// and perl6.vim (TODO give credit)
CodeMirror.defineMode("perl6",function(config,parserConfig){
	
	var PERL = {				
		// null - magic touch	
		//   1 - keyword
		//   2 - def
		//   3 - atom
		//   4 - operator
		//   5 - variable-2 (predefined)
		//   [x,y] - x=1,2,3; y=must be defined if x{...}
		'->'				:   4,
		'++'				:   4,
		'--'				:   4,
		'**'				:   4,
							//   ! ~ \ and unary + and -
		'=~'				:   4,
		'!~'				:   4,
		'*'				:   4,
		'/'				:   4,
		'%'				:   4,
		'x'				:   4,
		'+'				:   4,
		'-'				:   4,
		'.'				:   4,
		'<<'				:   4,
		'>>'				:   4,
							//   named unary operators
		'<'				:   4,
		'>'				:   4,
		'<='				:   4,
		'>='				:   4,
		'=='				:   4,
		'!='				:   4,
		'<=>'				:   4,
		'~~'				:   4,
		'&'				:   4,
		'|'				:   4,
		'^'				:   4,
		'&&'				:   4,
		'||'				:   4,
		'//'				:   4,
		'..'				:   4,
		'...'				:   4,
		'?'				:   4,
		':'				:   4,
		'='				:   4,
		'+='				:   4,
		'-='				:   4,
		'*='				:   4,	//   etc. ???
		','				:   4,
		'=>'				:   4,
		'::'				:   4,
		
		// PERL predefined variables (I know, what this is a paranoid idea, 
		// but may be needed for people, who learn PERL, and for me as well, 
		// ...and may be for you?;)
		'BEGIN'				:   [5,1],
		'END'				:   [5,1],
		'GETC'				:   [5,1],
		'READ'				:   [5,1],
		'READLINE'			:   [5,1],
		'DESTROY'			:   [5,1],
		'TIE'				:   [5,1],
		'TIEHANDLE'			:   [5,1],
		'UNTIE'				:   [5,1],
		'$ARG'				:    5,
		'$_'				:    5,
		'@ARG'				:    5,
		'@_'				:    5,
		'$"'				:    5,
		'$$'				:    5,
		'$('				:    5,
		'$EGID'				:    5,
		'$)'				:    5,
		'$0'				:    5,
		'$SUBSEP'			:    5,
		'$;'				:    5,
		'$UID'				:    5,
		'$<'				:    5,
		'$EUID'				:    5,
		'$>'				:    5,
		'$a'				:    5,
		'$b'				:    5,
		'$^C'				:    5,
		'$^D'				:    5,
		'$ENV'				:    5,
		'%ENV'				:    5,
		'$^F'				:    5,
		'@F'				:    5,
		'$^H'				:    5,
		'%^H'				:    5,
		'$INPLACE_EDIT'			:    5,
		'$^I'				:    5,
		'$^M'				:    5,
		'$OSNAME'			:    5,
		'$^O'				:    5,
		'${^OPEN}'			:    5,
		'$PERLDB'			:    5,
		'$^P'				:    5,
		'$SIG'				:    5,
		'%SIG'				:    5,
		'$BASETIME'			:    5,
		'$^T'				:    5,
		'${^TAINT}'			:    5,
		'${^UNICODE}'			:    5,
		'${^UTF8CACHE}'			:    5,
		'${^UTF8LOCALE}'		:    5,
		'$PERL_VERSION'			:    5,
		'$^V'				:    5,
		'${^WIN32_SLOPPY_STAT}'		:    5,
		'$EXECUTABLE_NAME'		:    5,
		'$^X'				:    5,
		'$1'				:    5,	// - regexp $1, $2...
		'$MATCH'			:    5,
		'$&'				:    5,
		'${^MATCH}'			:    5,
		'$PREMATCH'			:    5,
		'$`'				:    5,
		'${^PREMATCH}'			:    5,
		'$POSTMATCH'			:    5,
		"$'"				:    5,
		'${^POSTMATCH}'			:    5,
		'$LAST_PAREN_MATCH'		:    5,
		'$+'				:    5,
		'$LAST_SUBMATCH_RESULT'		:    5,
		'$^N'				:    5,
		'@LAST_MATCH_END'		:    5,
		'@+'				:    5,
		'%LAST_PAREN_MATCH'		:    5,
		'%+'				:    5,
		'@LAST_MATCH_START'		:    5,
		'@-'				:    5,
		'%LAST_MATCH_START'		:    5,
		'%-'				:    5,
		'$LAST_REGEXP_CODE_RESULT'	:    5,
		'$^R'				:    5,
		'${^RE_DEBUG_FLAGS}'		:    5,
		'${^RE_TRIE_MAXBUF}'		:    5,
		'$ARGV'				:    5,
		'@ARGV'				:    5,
		'ARGV'				:    5,
		'ARGVOUT'			:    5,
		'$OUTPUT_FIELD_SEPARATOR'	:    5,
		'$OFS'				:    5,
		'$,'				:    5,
		'$INPUT_LINE_NUMBER'		:    5,
		'$NR'				:    5,
		'$.'				:    5,
		'$INPUT_RECORD_SEPARATOR'	:    5,
		'$RS'				:    5,
		'$/'				:    5,
		'$OUTPUT_RECORD_SEPARATOR'	:    5,
		'$ORS'				:    5,
		'$\\'				:    5,
		'$OUTPUT_AUTOFLUSH'		:    5,
		'$|'				:    5,
		'$ACCUMULATOR'			:    5,
		'$^A'				:    5,
		'$FORMAT_FORMFEED'		:    5,
		'$^L'				:    5,
		'$FORMAT_PAGE_NUMBER'		:    5,
		'$%'				:    5,
		'$FORMAT_LINES_LEFT'		:    5,
		'$-'				:    5,
		'$FORMAT_LINE_BREAK_CHARACTERS'	:    5,
		'$:'				:    5,
		'$FORMAT_LINES_PER_PAGE'	:    5,
		'$='				:    5,
		'$FORMAT_TOP_NAME'		:    5,
		'$^'				:    5,
		'$FORMAT_NAME'			:    5,
		'$~'				:    5,
		'${^CHILD_ERROR_NATIVE}'	:    5,
		'$EXTENDED_OS_ERROR'		:    5,
		'$^E'				:    5,
		'$EXCEPTIONS_BEING_CAUGHT'	:    5,
		'$^S'				:    5,
		'$WARNING'			:    5,
		'$^W'				:    5,
		'${^WARNING_BITS}'		:    5,
		'$OS_ERROR'			:    5,
		'$ERRNO'			:    5,
		'$!'				:    5,
		'%OS_ERROR'			:    5,
		'%ERRNO'			:    5,
		'%!'				:    5,
		'$CHILD_ERROR'			:    5,
		'$?'				:    5,
		'$EVAL_ERROR'			:    5,
		'$@'				:    5,
		'$OFMT'				:    5,
		'$#'				:    5,
		'$*'				:    5,
		'$ARRAY_BASE'			:    5,
		'$['				:    5,
		'$OLD_PERL_VERSION'		:    5,
		'$]'				:    5,
						//	PERL blocks
		'if'				:[1,1],
		elsif				:[1,1],
		'else'				:[1,1],
		'while'				:[1,1],
		unless				:[1,1],
		'for'				:[1,1],
		foreach				:[1,1],
		q				:null,	// - singly quote a string
		qq				:null,	// - doubly quote a string
		qr				:null,	// - Compile pattern
		quotemeta			:null,	// - quote regular expression magic characters
		qw				:null,	// - quote a list of words
		qx				:null,	// - backquote quote a string
		y				:null};	// - transliterate a string


	var p6DeclareRoutine = [
 		"macro sub submethod method multi proto only rule token regex category",
 	];
	var p6Module = [
 		"module class role package enum grammar slang subset",
 	];
	var p6Variable = [
 		"self",
	];
	var p6Include = [
		"use require",
	];
	var p6Conditional = [
		"if else elsif unless",
	];
	var p6VarStorage = [
		"let my our state temp has constant",
	];
 	var p6Repeat = [
		"for loop repeat while until gather given",
	];
	var p6FlowControl = [
		"take do when next last redo return contend maybe defer",
		"default exit make continue break goto leave async lift",
	];
	var p6TypeConstraint = [
		"is as but trusts of returns handles where augment supersede",
	];
 	var p6ClosureTrait = [
		"BEGIN CHECK INIT START FIRST ENTER LEAVE KEEP",
		"UNDO NEXT LAST PRE POST END CATCH CONTROL TEMP",
	];
	var p6Exception = [
		"die fail try warn",
	];
 	var p6Property = [
		"prec irs ofs ors export deep binary unary reparsed rw parsed cached",
		"readonly defequiv will ref copy inline tighter looser equiv assoc",
		"required",
	];
	var p6Number = [
		"NaN Inf",
	];
	var p6Pragma = [
		"oo fatal",	
	];

	var p6Type = [
		'Object Any Junction Whatever Capture Match',
		'Signature Proxy Matcher Package Module Class',
		'Grammar Scalar Array Hash KeyHash KeySet KeyBag',
        'Pair List Seq Range Set Bag Mapping Void Undef',
        'Failure Exception Code Block Routine Sub Macro',
        'Method Submethod Regex Str Blob Char Byte',
        'Codepoint Grapheme StrPos StrLen Version Num',
        'Complex num complex Bit bit bool True False',
        'Increasing Decreasing Ordered Callable AnyChar',
		'Positional Associative Ordering KeyExtractor',
		'Comparator OrderingPair IO KitchenSink Role',
		'Int int int1 int2 int4 int8 int16 int32 int64',
		'Rat rat rat1 rat2 rat4 rat8 rat16 rat32 rat64',
		'Buf buf buf1 buf2 buf4 buf8 buf16 buf32 buf64',
		'UInt uint uint1 uint2 uint4 uint8 uint16 uint32',
		'uint64 Abstraction utf8 utf16 utf32'
	];
	var p6Routines = [
    	"eager hyper substr index rindex grep map sort join lines hints chmod",
		"split reduce min max reverse truncate zip cat roundrobin classify",
    	"first sum keys values pairs defined delete exists elems end kv any",
    	"all one wrap shape key value name pop push shift splice unshift floor",
    	"ceiling abs exp log log10 rand sign sqrt sin cos tan round strand",
    	"roots cis unpolar polar atan2 pick chop p5chop chomp p5chomp lc",
    	"lcfirst uc ucfirst capitalize normalize pack unpack quotemeta comb",
    	"samecase sameaccent chars nfd nfc nfkd nfkc printf sprintf caller",
  	  	"evalfile run runinstead nothing want bless chr ord gmtime time eof",
    	"localtime gethost getpw chroot getlogin getpeername kill fork wait",
    	"perl graphs codes bytes clone print open read write readline say seek",
    	"close opendir readdir slurp pos fmt vec link unlink symlink uniq pair",
    	"asin atan sec cosec cotan asec acosec acotan sinh cosh tanh asinh",
    	"acos acosh atanh sech cosech cotanh sech acosech acotanh asech ok",
    	"plan_ok dies_ok lives_ok skip todo pass flunk force_todo use_ok isa_ok",
    	"diag is_deeply isnt like skip_rest unlike cmp_ok eval_dies_ok nok_error",
    	"eval_lives_ok approx is_approx throws_ok version_lt plan eval succ pred",
    	"times nonce once signature new connect operator undef undefine sleep",
    	"from to infix postfix prefix circumfix postcircumfix minmax lazy count",
    	"unwrap getc pi e context void quasi body each contains rewinddir subst",
    	"can isa flush arity assuming rewind callwith callsame nextwith nextsame",
    	"attr eval_elsewhere none srand trim trim_start trim_end lastcall WHAT",
    	"WHERE HOW WHICH VAR WHO WHENCE ACCEPTS REJECTS does not true iterator by",
    	"re im invert flip",
	];
	var p6Operator = [
		"div x xx mod also leg cmp before after eq ne le lt",
		"gt ge eqv ff fff and andthen Z X or xor",
		"orelse extra m mm rx s tr",
	];

	var addWords = function(p6Words, category) {
		var i, j, words;
		for (i = 0; i < p6Words.length; i++) {
        	words = p6Words[i].split(' ');
        	for(j = 0; j < words.length; j++) {
            	PERL[words[j]] = category;
        	}
    	}
	};

	addWords(p6DeclareRoutine, 1);
	addWords(p6Module,         1);
    addWords(p6Variable,       1);
    addWords(p6Include,        1);
    addWords(p6Conditional,    1);
    addWords(p6VarStorage,     2);
    addWords(p6Repeat,         1);
    addWords(p6FlowControl,    1);
    addWords(p6TypeConstraint, 1);
    addWords(p6ClosureTrait,   1);
    addWords(p6Exception,      1);
    addWords(p6Property,       1);
    addWords(p6Number,         1);
    addWords(p6Pragma,         1);
	addWords(p6Type,           2);
	addWords(p6Routines,       3);
	addWords(p6Operator,       4);

	var RXstyle="string-2";
	var RXmodifiers=/[goseximacplud]/;		// NOTE: "m", "s", "y" and "tr" need to correct real modifiers for each regexp type

	function tokenChain(stream,state,chain,style,tail){	// NOTE: chain.length > 2 is not working now (it's for s[...][...]geos;)
		state.chain=null;                               //                                                          12   3tail
		state.style=null;
		state.tail=null;
		state.tokenize=function(stream,state){
			var e=false,c,i=0;
			while(c=stream.next()){
				if(c===chain[i]&&!e){
					if(chain[++i]!==undefined){
						state.chain=chain[i];
						state.style=style;
						state.tail=tail;}
					else if(tail)
						stream.eatWhile(tail);
					state.tokenize=tokenPerl;
					return style;}
				e=!e&&c=="\\";}
			return style;};
		return state.tokenize(stream,state);}

	function tokenSOMETHING(stream,state,string){
		state.tokenize=function(stream,state){
			if(stream.string==string)
				state.tokenize=tokenPerl;
			stream.skipToEnd();
			return "string";};
		return state.tokenize(stream,state);}

	function tokenPerl(stream,state){
		if(stream.eatSpace())
			return null;
		if(state.chain)
			return tokenChain(stream,state,state.chain,state.style,state.tail);
		if(stream.match(/^\-?[\d\.]/,false))
			if(stream.match(/^(\-?(\d*\.\d+(e[+-]?\d+)?|\d+\.\d*)|0x[\da-fA-F]+|0b[01]+|\d+(e[+-]?\d+)?)/))
				return 'number';
		if(stream.match(/^<<(?=\w)/)){			// NOTE: <<SOMETHING\n...\nSOMETHING\n
			stream.eatWhile(/\w/);
			return tokenSOMETHING(stream,state,stream.current().substr(2));}
		if(stream.sol()&&stream.match(/^\=item(?!\w)/)){// NOTE: \n=item...\n=cut\n
			return tokenSOMETHING(stream,state,'=cut');}
		var ch=stream.next();
		if(ch=='"'||ch=="'"){				// NOTE: ' or " or <<'SOMETHING'\n...\nSOMETHING\n or <<"SOMETHING"\n...\nSOMETHING\n
			if(stream.prefix(3)=="<<"+ch){
				var p=stream.pos;
				stream.eatWhile(/\w/);
				var n=stream.current().substr(1);
				if(n&&stream.eat(ch))
					return tokenSOMETHING(stream,state,n);
				stream.pos=p;}
			return tokenChain(stream,state,[ch],"string");}
		if(ch=="q"){
			var c=stream.look(-2);
			if(!(c&&/\w/.test(c))){
				c=stream.look(0);
				if(c=="x"){
					c=stream.look(1);
					if(c=="("){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[")"],RXstyle,RXmodifiers);}
					if(c=="["){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["]"],RXstyle,RXmodifiers);}
					if(c=="{"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["}"],RXstyle,RXmodifiers);}
					if(c=="<"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[">"],RXstyle,RXmodifiers);}
					if(/[\^'"!~\/]/.test(c)){
						stream.eatSuffix(1);
						return tokenChain(stream,state,[stream.eat(c)],RXstyle,RXmodifiers);}}
				else if(c=="q"){
					c=stream.look(1);
					if(c=="("){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[")"],"string");}
					if(c=="["){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["]"],"string");}
					if(c=="{"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["}"],"string");}
					if(c=="<"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[">"],"string");}
					if(/[\^'"!~\/]/.test(c)){
						stream.eatSuffix(1);
						return tokenChain(stream,state,[stream.eat(c)],"string");}}
				else if(c=="w"){
					c=stream.look(1);
					if(c=="("){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[")"],"bracket");}
					if(c=="["){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["]"],"bracket");}
					if(c=="{"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["}"],"bracket");}
					if(c=="<"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[">"],"bracket");}
					if(/[\^'"!~\/]/.test(c)){
						stream.eatSuffix(1);
						return tokenChain(stream,state,[stream.eat(c)],"bracket");}}
				else if(c=="r"){
					c=stream.look(1);
					if(c=="("){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[")"],RXstyle,RXmodifiers);}
					if(c=="["){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["]"],RXstyle,RXmodifiers);}
					if(c=="{"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,["}"],RXstyle,RXmodifiers);}
					if(c=="<"){
						stream.eatSuffix(2);
						return tokenChain(stream,state,[">"],RXstyle,RXmodifiers);}
					if(/[\^'"!~\/]/.test(c)){
						stream.eatSuffix(1);
						return tokenChain(stream,state,[stream.eat(c)],RXstyle,RXmodifiers);}}
				else if(/[\^'"!~\/(\[{<]/.test(c)){
					if(c=="("){
						stream.eatSuffix(1);
						return tokenChain(stream,state,[")"],"string");}
					if(c=="["){
						stream.eatSuffix(1);
						return tokenChain(stream,state,["]"],"string");}
					if(c=="{"){
						stream.eatSuffix(1);
						return tokenChain(stream,state,["}"],"string");}
					if(c=="<"){
						stream.eatSuffix(1);
						return tokenChain(stream,state,[">"],"string");}
					if(/[\^'"!~\/]/.test(c)){
						return tokenChain(stream,state,[stream.eat(c)],"string");}}}}
		if(ch=="m"){
			var c=stream.look(-2);
			if(!(c&&/\w/.test(c))){
				c=stream.eat(/[(\[{<\^'"!~\/]/);
				if(c){
					if(/[\^'"!~\/]/.test(c)){
						return tokenChain(stream,state,[c],RXstyle,RXmodifiers);}
					if(c=="("){
						return tokenChain(stream,state,[")"],RXstyle,RXmodifiers);}
					if(c=="["){
						return tokenChain(stream,state,["]"],RXstyle,RXmodifiers);}
					if(c=="{"){
						return tokenChain(stream,state,["}"],RXstyle,RXmodifiers);}
					if(c=="<"){
						return tokenChain(stream,state,[">"],RXstyle,RXmodifiers);}}}}
		if(ch=="s"){
			var c=/[\/>\]})\w]/.test(stream.look(-2));
			if(!c){
				c=stream.eat(/[(\[{<\^'"!~\/]/);
				if(c){
					if(c=="[")
						return tokenChain(stream,state,["]","]"],RXstyle,RXmodifiers);
					if(c=="{")
						return tokenChain(stream,state,["}","}"],RXstyle,RXmodifiers);
					if(c=="<")
						return tokenChain(stream,state,[">",">"],RXstyle,RXmodifiers);
					if(c=="(")
						return tokenChain(stream,state,[")",")"],RXstyle,RXmodifiers);
					return tokenChain(stream,state,[c,c],RXstyle,RXmodifiers);}}}
		if(ch=="y"){
			var c=/[\/>\]})\w]/.test(stream.look(-2));
			if(!c){
				c=stream.eat(/[(\[{<\^'"!~\/]/);
				if(c){
					if(c=="[")
						return tokenChain(stream,state,["]","]"],RXstyle,RXmodifiers);
					if(c=="{")
						return tokenChain(stream,state,["}","}"],RXstyle,RXmodifiers);
					if(c=="<")
						return tokenChain(stream,state,[">",">"],RXstyle,RXmodifiers);
					if(c=="(")
						return tokenChain(stream,state,[")",")"],RXstyle,RXmodifiers);
					return tokenChain(stream,state,[c,c],RXstyle,RXmodifiers);}}}
		if(ch=="t"){
			var c=/[\/>\]})\w]/.test(stream.look(-2));
			if(!c){
				c=stream.eat("r");if(c){
				c=stream.eat(/[(\[{<\^'"!~\/]/);
				if(c){
					if(c=="[")
						return tokenChain(stream,state,["]","]"],RXstyle,RXmodifiers);
					if(c=="{")
						return tokenChain(stream,state,["}","}"],RXstyle,RXmodifiers);
					if(c=="<")
						return tokenChain(stream,state,[">",">"],RXstyle,RXmodifiers);
					if(c=="(")
						return tokenChain(stream,state,[")",")"],RXstyle,RXmodifiers);
					return tokenChain(stream,state,[c,c],RXstyle,RXmodifiers);}}}}
		if(ch=="`"){
			return tokenChain(stream,state,[ch],"variable-2");}
		if(ch=="/"){
			if(!/~\s*$/.test(stream.prefix()))
				return "operator";
			else
				return tokenChain(stream,state,[ch],RXstyle,RXmodifiers);}
		if(ch=="$"){
			var p=stream.pos;
			if(stream.eatWhile(/\d/)||stream.eat("{")&&stream.eatWhile(/\d/)&&stream.eat("}"))
				return "variable-2";
			else
				stream.pos=p;}
		if(/[$@%]/.test(ch)){
			var p=stream.pos;
			if(stream.eat("^")&&stream.eat(/[A-Z]/)||!/[@$%&]/.test(stream.look(-2))&&stream.eat(/[=|\\\-#?@;:&`~\^!\[\]*'"$+.,\/<>()]/)){
				var c=stream.current();
				if(PERL[c])
					return "variable-2";}
			stream.pos=p;}
		if(/[$@%&]/.test(ch)){
			if(stream.eatWhile(/[\w$\[\]]/)||stream.eat("{")&&stream.eatWhile(/[\w$\[\]]/)&&stream.eat("}")){
				var c=stream.current();
				if(PERL[c])
					return "variable-2";
				else
					return "variable";}}
		if(ch=="#"){
			if(stream.look(-2)!="$"){
				stream.skipToEnd();
				return "comment";}}
		if(/[:+\-\^*$&%@=<>!?|\/~\.]/.test(ch)){
			var p=stream.pos;
			stream.eatWhile(/[:+\-\^*$&%@=<>!?|\/~\.]/);
			if(PERL[stream.current()])
				return "operator";
			else
				stream.pos=p;}
		if(ch=="_"){
			if(stream.pos==1){
				if(stream.suffix(6)=="_END__"){
					return tokenChain(stream,state,['\0'],"comment");}
				else if(stream.suffix(7)=="_DATA__"){
					return tokenChain(stream,state,['\0'],"variable-2");}
				else if(stream.suffix(7)=="_C__"){
					return tokenChain(stream,state,['\0'],"string");}}}
		if(/\w/.test(ch)){
			var p=stream.pos;
			if(stream.look(-2)=="{"&&(stream.look(0)=="}"||stream.eatWhile(/\w/)&&stream.look(0)=="}"))
				return "string";
			else
				stream.pos=p;}
		if(/[A-Z]/.test(ch)){
			var l=stream.look(-2);
			var p=stream.pos;
			stream.eatWhile(/[A-Z_]/);
			if(/[\da-z]/.test(stream.look(0))){
				stream.pos=p;}
			else{
				var c=PERL[stream.current()];
				if(!c)
					return "meta";
				if(c[1])
					c=c[0];
				if(l!=":"){
					if(c==1)
						return "keyword";
					else if(c==2)
						return "def";
					else if(c==3)
						return "atom";
					else if(c==4)
						return "operator";
					else if(c==5)
						return "variable-2";
					else
						return "meta";}
				else
					return "meta";}}
		if(/[a-zA-Z_]/.test(ch)){
			var l=stream.look(-2);
			stream.eatWhile(/\w/);
			var c=PERL[stream.current()];
			if(!c)
				return "meta";
			if(c[1])
				c=c[0];
			if(l!=":"){
				if(c==1)
					return "keyword";
				else if(c==2)
					return "def";
				else if(c==3)
					return "atom";
				else if(c==4)
					return "operator";
				else if(c==5)
					return "variable-2";
				else
					return "meta";}
			else
				return "meta";}
		return null;}

	return{
		startState:function(){
			return{
				tokenize:tokenPerl,
				chain:null,
				style:null,
				tail:null};},
		token:function(stream,state){
			return (state.tokenize||tokenPerl)(stream,state);},
		electricChars:"{}"};});

CodeMirror.defineMIME("text/x-perl6", "perl6");

// it's like "peek", but need for look-ahead or look-behind if index < 0
CodeMirror.StringStream.prototype.look=function(c){
	return this.string.charAt(this.pos+(c||0));};

// return a part of prefix of current stream from current position
CodeMirror.StringStream.prototype.prefix=function(c){
	if(c){
		var x=this.pos-c;
		return this.string.substr((x>=0?x:0),c);}
	else{
		return this.string.substr(0,this.pos-1);}};

// return a part of suffix of current stream from current position
CodeMirror.StringStream.prototype.suffix=function(c){
	var y=this.string.length;
	var x=y-this.pos+1;
	return this.string.substr(this.pos,(c&&c<y?c:x));};

// return a part of suffix of current stream from current position and change current position
CodeMirror.StringStream.prototype.nsuffix=function(c){
	var p=this.pos;
	var l=c||(this.string.length-this.pos+1);
	this.pos+=l;
	return this.string.substr(p,l);};

// eating and vomiting a part of stream from current position
CodeMirror.StringStream.prototype.eatSuffix=function(c){
	var x=this.pos+c;
	var y;
	if(x<=0)
		this.pos=0;
	else if(x>=(y=this.string.length-1))
		this.pos=y;
	else
		this.pos=x;};
