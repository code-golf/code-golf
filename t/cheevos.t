use t;

constant %holes = <
    1  {hello-world}
    11 {up-to-eleven}
    13 {bakers-dozen}
    19 {the-watering-hole}
    21 {blackjack}
    34 {rule-34}
    40 {forty-winks}
    42 {dont-panic}
    50 {bullseye}
    60 {gone-in-60-holes}
    69 {cunning-linguist}
>;

constant %langs = <
    12 {polyglot}
    24 {polyglutton}
>;

my $dbh = dbh;
createUser($dbh, 1);
my $session = createSession($dbh, 1);

$client.get: 'https://app:1443/about',
    headers => { cookie => "__Host-session=$session" };

is $dbh.execute('SELECT ARRAY(SELECT trophy FROM trophies)').row, '{rtfm}',
    'GET /about earns {rtfm}';

for $dbh.execute('SELECT unnest(enum_range(null::hole))').allrows.flat {
    my $cheevos = %holes{ my $i = ++$ } // '{}';

    # Add hole-specific cheevos on the end.
    $cheevos.=subst: '}', ',interview-ready}' if $_ eq 'fizz-buzz';
    $cheevos.=subst: '}', ',solve-quine}'     if $_ eq 'quine';
    $cheevos.=subst: '{,', '{';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', ?, 'c', 1)", $_,
    ).row, $cheevos, "$_/c earns $cheevos";
}

for <
    brainfuck       brainfuck {inception}
    divisors        php       {elephpant-in-the-room}
    poker           fish      {fish-n-chips}
    quine           python    {ouroboros}
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

for $dbh.execute('SELECT unnest(enum_range(null::lang))').allrows.flat {
    next if $_ eq 'assembly';   # Assembly can't have non-NULL chars.

    my $earns = %langs{ my $i = ++$ } // '{}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', 'λ', ?, 1)", $_,
    ).row, $earns, "Lang $i earns $earns";
}

done-testing;
