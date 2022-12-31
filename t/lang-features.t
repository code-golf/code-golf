use t;

is post-solution(|.value)<Err>, '', .key for
    awk-ordchr  => \(:lang<awk>    :code('@load "ordchr"')),
    nim-re      => \(:lang<nim>    :code('import re;echo "a".match(re"a")')),
    perl-bigint => \(:lang<perl>   :code('use bigint')),
    perl-bignum => \(:lang<perl>   :code('use bignum')),
    perl-glob   => \(:lang<perl>   :code('<foo{bar,baz}>')),
    prolog-re   => \(:lang<prolog> :code(':- crypto_is_prime(5, []).')),
    raku-exp    => \(:lang<raku>   :code('use experimental')),
    tcl-min     => \(:lang<tcl>    :code('puts [expr min(5,6)]'));

done-testing;
