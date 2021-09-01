use t;
use JSON::Fast;

my $dbh = dbh;
createUser($dbh, 1);
my $session = createSession($dbh, 1);

is-deeply export, {
    name      => 'Bob',
    country   => Nil,
    time_zone => Nil,
    cheevos   => [],
    solutions => [],
}, 'Initial export is empty';

is-deeply export, {
    name      => 'Bob',
    country   => Nil,
    time_zone => Nil,
    cheevos   => [ { :cheevo<takeout> :earned<datetime> }, ],
    solutions => [],
}, 'Second export has takeout cheevo';

my $code = 'say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1 .. 100';

post-solution :$code :$session;

is-deeply export, {
    name      => 'Bob',
    country   => Nil,
    time_zone => Nil,
    cheevos   => [
        { :cheevo<hello-world>     :earned<datetime> },
        { :cheevo<interview-ready> :earned<datetime> },
        { :cheevo<takeout>         :earned<datetime> },
    ],
    solutions => [
        {
            :hole<fizz-buzz>, :lang<raku>, :scoring<bytes>, :!failing,
            :bytes(58), :chars(58), :submitted<datetime>, :$code,
        },
        {
            :hole<fizz-buzz>, :lang<raku>, :scoring<chars>, :!failing,
            :bytes(58), :chars(58), :submitted<datetime>, :$code,
        },
    ],
}, 'Third export has multiple cheevos & solutions';

sub export {
    my %res = $client.get: 'https://app:1443/golfer/export',
        headers => { cookie => "__Host-session=$session" };

    once {
        is %res<headers><content-disposition>,
            'attachment; filename="code-golf-export.json"',
            'Content-Disposition header';
    }

    # Replace all valid datetimes with a string for is-deeply matching later.
    return from-json %res<content>.decode.subst:
        / \d**4 '-' \d\d '-' \d\d 'T' \d\d ':' \d\d ':' \d\d '.' \d+ 'Z' /,
        'datetime', :g;
}

done-testing;
