use t;

is post-solution(|.value)<runs>[0]<stderr>, '', .key for
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
like post-solution(:lang<j> :code('echo JVERSION'))<runs>[0]<stdout>,
    / '/j64/linux' /, 'J engine is baseline AMD 64 (no AVX 512)';

# Null byte in solution.
is-deeply post-solution( :code(qq:to/CODE/) )<runs>[0]<pass stderr>:p,
        say "Fizz" x \$_ %% 3 ~ "Buzz" x \$_ %% 5 || \$_ for 1 .. 100;
        # Null byte in comment = \x00
    CODE
    ( :!pass :stderr('Solutions must not contain a literal null byte.') ),
    'Null byte in solution fails';

# Trivial Tex Quine.
my $code = "Trivial\n";
my $err  = ｢Quine in TeX must have at least one '\' character.｣;
my %res  = post-solution :hole<quine> :lang<tex> :$code;

is-deeply %res<runs>[0]<answer pass stderr>:p,
    ( :answer($code) :!pass :stderr($err) ), 'Trivial Tex Quine is blocked';

done-testing;
