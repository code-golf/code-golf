use t;

is post-solution(|.value)<ExitCode>, 0, .key for
    js-i18n   => \(:lang<javascript> :code('Intl')),
    nim-re    => \(:lang<nim>        :code('import re;echo "a".match(re"a")')),
    perl-glob => \(:lang<perl>       :code('<foo{bar,baz}>'));

done-testing;
