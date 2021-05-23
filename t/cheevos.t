use t;

constant %cheevos = <
    1  {hello-world}
    11 {up-to-eleven}
    13 {bakers-dozen}
    19 {the-watering-hole}
    40 {forty-winks}
    42 {dont-panic}
    50 {bullseye}
>;

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (1, '')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (1) RETURNING id').row.head;

await $client.get: 'https://app:1443/about',
    headers => [ cookie => "__Host-session=$session" ];

is $dbh.execute('SELECT ARRAY(SELECT trophy FROM trophies)').row, '{rtfm}',
    'GET /about earns {rtfm}';

for from-toml slurp 'holes.toml' {
    next if .value<experiment>;     # Experimental holes can't be saved.
    next if .key eq 'Fizz Buzz';    # This is tested lower.

    my $hole    = .key.lc.trans: ' ’' => '-', :d;
    my $cheevos = %cheevos{ my $i = ++$ } // '{}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', ?, 'c', 1)", $hole,
    ).row, $cheevos, "Solution $i earns $cheevos";
}

for <
    brainfuck       brainfuck {inception}
    divisors        php       {elephpant-in-the-room}
    fizz-buzz       haskell   {interview-ready}
    poker           fish      {fish-n-chips}
    quine           python    {ouroboros}
    ten-pin-bowling cobol     {cobowl}
> -> $hole, $lang, $cheevos {
    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', ?, ?, 1)", $hole, $lang,
    ).row, $cheevos, "$hole/$lang earns $cheevos";
}

is $dbh.execute(
    "SELECT earned FROM save_solution(3, 1, '⛳', 'π', 'c', 1)",
).row, '{different-strokes}', 'Earns {different-strokes}';

done-testing;
