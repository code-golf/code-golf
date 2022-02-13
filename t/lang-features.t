use t;

is post-solution(|.value)<Err>, '', .key for
    nim-re    => \(:lang<nim>    :code('import re;echo "a".match(re"a")')),
    perl-glob => \(:lang<perl>   :code('<foo{bar,baz}>')),
    prolog-re => \(:lang<prolog> :code(':- crypto_is_prime(5, []).')),
    raku-exp  => \(:lang<raku>   :code('use experimental'));

done-testing;
