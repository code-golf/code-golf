use t;

constant $code-long  = 'say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1 .. 100';
constant $code-short = 'say "Fizz"x$_%%3~"Buzz"x$_%%5||$_ for 1..100';

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'test')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

subtest 'Failing solution' => {
    nok post-solution( :$session, :code('say 1') )<Pass>, 'Solution fails';

    is $dbh.execute('SELECT COUNT(*) FROM code').row.head, 0, 'DB is empty';
}

subtest 'Initial solution' => {
    ok post-solution( :$session, :code($code-long) )<Pass>, 'Solution passes';

    is-deeply db(), (
        { code => $code-long, code_id => 1, user_id => 123 },
    ), 'DB is inserted';
};

subtest 'Same solution' => {
    ok post-solution( :$session, :code($code-long) )<Pass>, 'Solution passes';

    is-deeply db(), (
        { code => $code-long, code_id => 1, user_id => 123 },
    ), 'DB is the same';
};

subtest 'Updated solution' => {
    ok post-solution( :$session, :code($code-short) )<Pass>, 'Solution passes';

    # Note how the sequence has jumped because of the second solution.
    is-deeply db(), (
        { code => $code-short, code_id => 3, user_id => 123 },
    ), 'DB is updated';

    todo "Code isn't cleaned up yet";

    is $dbh.execute('SELECT COUNT(*) FROM code').row.head, 1, 'Code is cleaned up';
};

subtest 'Different user' => {
    $dbh.execute: "INSERT INTO users (id, login) VALUES (456, 'test2')";

    my $session = $dbh.execute(
        'INSERT INTO sessions (user_id) VALUES (456) RETURNING id').row.head;

    ok post-solution( :$session, :code($code-short) )<Pass>, 'Solution passes';

    # Note how they share the some code ID.
    is-deeply db(), (
        { code => $code-short, code_id => 3, user_id => 123 },
        { code => $code-short, code_id => 3, user_id => 456 },
    ), 'DB is updated';
};

sub db {
    $dbh.execute(
        'SELECT code, code_id, user_id FROM solutions JOIN code ON id = code_id',
    ).allrows(:array-of-hash)
}

done-testing;
