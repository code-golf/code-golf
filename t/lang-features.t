use t;

is post-solution(|.value)<Err>, '', .key for
    awk-ordchr  => \(:lang<awk>        :code('@load "ordchr"')),
    js-icu      => \(:lang<javascript> :code('/\p{Emoji}/u')),
    nim-re      => \(:lang<nim>        :code('import re;echo "a".match(re"a")')),
    perl-bigint => \(:lang<perl>       :code('use bigint')),
    perl-bignum => \(:lang<perl>       :code('use bignum')),
    perl-glob   => \(:lang<perl>       :code('<foo{bar,baz}>')),
    prolog-re   => \(:lang<prolog>     :code(':- crypto_is_prime(5, []).')),
    raku-exp    => \(:lang<raku>       :code('use experimental')),
    tcl-min     => \(:lang<tcl>        :code('puts [expr min(5,6)]'));

# AVX 512 wouldn't work on live, yet.
like post-solution(:lang<j> :code('echo JVERSION'))<Out>, / '/j64/linux' /,
    'J engine is baseline AMD 64 (no AVX 512)';

# Trivial Tex Quine.
my $code = "Trivial\n";
my $err  = 'Quine in TeX must have at least one &#39;\&#39; character.';
my %res  = post-solution :hole<quine> :lang<tex> :$code;

is-deeply %res<Err Exp Pass>:p, ( :Err($err) :Exp($code) :!Pass ),
    'Trivial Tex Quine is blocked';

done-testing;
