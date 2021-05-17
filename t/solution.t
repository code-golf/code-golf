use t;

constant $code-long        = 'say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1 .. 100';
constant $code-short       = 'say "Fizz"x$_%%3~"Buzz"x$_%%5||$_ for 1..100';
constant $code-short-chars = 'say "Fizz"x$_%%3~"Buzz"x$_%%5||$_ for 1…100';

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'test')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

subtest 'Failing solution' => {
    nok post-solution( :$session, :code('say 1') )<Pass>, 'Solution fails';

    is $dbh.execute('SELECT COUNT(*) FROM solutions').row.head, 0, 'DB is empty';
}

subtest 'Initial solution' => {
    ok post-solution( :$session, :code($code-long) )<Pass>, 'Passes';

    is-deeply db, (
        { code => $code-long, scoring => 'bytes' },
        { code => $code-long, scoring => 'chars' },
    ), 'Inserts both';
};

subtest 'Same solution' => {
    ok post-solution( :$session, :code($code-long) )<Pass>, 'Passes';

    is-deeply db, (
        { code => $code-long, scoring => 'bytes' },
        { code => $code-long, scoring => 'chars' },
    ), 'Updates none';
};

subtest 'Shorter solution' => {
    ok post-solution( :$session, :code($code-short) )<Pass>, 'Passes';

    is-deeply db, (
        { code => $code-short, scoring => 'bytes' },
        { code => $code-short, scoring => 'chars' },
    ), 'Updates both';
};

subtest 'Shorter chars solution' => {
    ok post-solution( :$session, :code($code-short-chars) )<Pass>, 'Passes';

    is-deeply db, (
        { code => $code-short,       scoring => 'bytes' },
        { code => $code-short-chars, scoring => 'chars' },
    ), 'Updates just the chars';
};

sub db {
    $dbh.execute('SELECT code, scoring FROM solutions').allrows :array-of-hash;
}

done-testing;
