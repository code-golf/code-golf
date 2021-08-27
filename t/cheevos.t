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
createUser($dbh, 1);
my $session = createSession($dbh, 1);

$client.get: 'https://app:1443/about',
    headers => { cookie => "__Host-session=$session" };

is $dbh.execute('SELECT ARRAY(SELECT trophy FROM trophies)').row, '{rtfm}',
    'GET /about earns {rtfm}';

for %( from-toml 'holes.toml'.IO ) {
    next if .value<experiment>;             # Experimental holes can't be saved.
    next if .key ~~ 'Fizz Buzz' | 'Quine';  # Theese are tested lower.

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
    quine           python    {solve-quine,ouroboros}
    seven-segment   assembly  {assembly-required}
    sudoku          hexagony  {off-the-grid}
    ten-pin-bowling cobol     {cobowl}
> -> $hole, $lang, $cheevos {
    is $dbh.execute(
        "SELECT earned FROM save_solution(2, ?, 'ab', ?, ?, 1)",
        $lang eq 'assembly' ?? Nil !! 2, $hole, $lang,
    ).row, $cheevos, "$hole/$lang earns $cheevos";
}

is $dbh.execute(
    "SELECT earned FROM save_solution(3, 1, '⛳', 'π', 'c', 1)",
).row, '{different-strokes}', 'Earns {different-strokes}';

done-testing;
