use t;

constant $answer           = slurp("config/hole-answers/fizz-buzz.txt").chomp;
constant $code-long        = 'say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1 .. 100';
constant $code-short       = 'say "Fizz"x$_%%3~"Buzz"x$_%%5||$_ for 1..100';
constant $code-short-chars = 'say "Fizz"x$_%%3~"Buzz"x$_%%5||$_ for 1…100';

is-deeply post-solution(:code($code-long))<runs>[0]<
    answer args exit_code pass stderr
>:kv.Hash, {
    answer    => $answer,
    args      => [],
    exit_code => 0,
    pass      => True,
    stderr    => '',
};

my $dbh     = dbh;
my $session = new-golfer :$dbh;

subtest 'Failing solution' => {
    nok post-solution( :$session :code('say 1') )<runs>[0]<pass>, 'Solution fails';

    is $dbh.execute('SELECT COUNT(*) FROM solutions').row.head, 0, 'DB is empty';
}

subtest 'Initial solution' => {
    ok post-solution( :$session :code($code-long) )<runs>[0]<pass>, 'Passes';

    is-deeply db, (
        { :code($code-long), :lang<raku>, :scoring<bytes> },
        { :code($code-long), :lang<raku>, :scoring<chars> },
    ), 'Inserts both';
};

subtest 'Same solution' => {
    ok post-solution( :$session :code($code-long) )<runs>[0]<pass>, 'Passes';

    is-deeply db, (
        { :code($code-long), :lang<raku>, :scoring<bytes> },
        { :code($code-long), :lang<raku>, :scoring<chars> },
    ), 'Updates none';
};

subtest 'Shorter solution' => {
    ok post-solution( :$session :code($code-short) )<runs>[0]<pass>, 'Passes';

    is-deeply db, (
        { :code($code-short), :lang<raku>, :scoring<bytes> },
        { :code($code-short), :lang<raku>, :scoring<chars> },
    ), 'Updates both';
};

subtest 'Shorter chars solution' => {
    ok post-solution( :$session :code($code-short-chars) )<runs>[0]<pass>, 'Passes';

    is-deeply db, (
        { :code($code-short),       :lang<raku>, :scoring<bytes> },
        { :code($code-short-chars), :lang<raku>, :scoring<chars> },
    ), 'Updates just the chars';
};

subtest 'Assembly solution' => {
constant $code = Q:c:to/ASM/;
    mov $1,              %eax
    mov $1,              %edi
    mov $text,           %rsi
    mov $textEnd - text, %edx
    syscall

    text: .string "{ $answer.lines.join: ｢\n｣ }"; textEnd:
    ASM

    ok post-solution( :$session :$code :lang<assembly> )<runs>[0]<pass>, 'Passes';

    is-deeply db, (
        { :code($code-short),       :lang<raku>,     :scoring<bytes> },
        { :code($code-short-chars), :lang<raku>,     :scoring<chars> },
        { :$code,                   :lang<assembly>, :scoring<bytes> },
    ), 'Inserts only bytes';
};

sub db {
    $dbh.execute(
        'SELECT code, lang, scoring FROM solutions').allrows :array-of-hash;
}

done-testing;
