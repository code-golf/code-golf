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
    80 {phileas-fogg}
    86 {x86}
>;

constant %langs = <
    12 {polyglot}
    24 {polyglutton}
    36 {omniglot}
>;

my $dbh = dbh;
createUser($dbh, 1);
my $session = createSession($dbh, 1);

$client.get: 'https://app/about',
    headers => { cookie => "__Host-session=$session" };

is $dbh.execute('SELECT ARRAY(SELECT trophy FROM trophies)').row, '{rtfm}',
    'GET /about earns {rtfm}';

for $dbh.execute('SELECT unnest(enum_range(null::hole))').allrows.flat {
    my $cheevos = %holes{ my $i = ++$ } // '{}';

    # Add hole-specific cheevos to the start.
    $cheevos.=subst: '{', '{interview-ready,' if $_ eq 'fizz-buzz';
    $cheevos.=subst: '{', '{solve-quine,'     if $_ eq 'quine';
    $cheevos.=subst: ',}', '}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', ?, 'c', 1)", $_,
    ).row, $cheevos, "$_/c earns $cheevos";
}

for <
    brainfuck        brainfuck {inception}
    divisors         php       {elephpant-in-the-room}
    hexdump          hexagony  {hextreme-agony}
    pascals-triangle pascal    {under-pressure}
    poker            fish      {fish-n-chips}
    quine            python    {ouroboros}
    seven-segment    assembly  {assembly-required}
    sudoku           hexagony  {off-the-grid}
    ten-pin-bowling  cobol     {cobowl}
> -> $hole, $lang, $cheevos {
    is $dbh.execute(
        "SELECT earned FROM save_solution(2, ?, 'ab', ?, ?, 1)",
        $lang eq 'assembly' ?? Nil !! 2, $hole, $lang,
    ).row, $cheevos, "$hole/$lang earns $cheevos";
}

for <brainfuck d hexagony javascript nim swift sql zig> {
    my $cheevos = $_ eq 'zig' ?? '{pangramglot}' !! '{}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', 'pangram-grep', ?, 1)",
        $_,
    ).row, $cheevos, "pangram-grep/$_ earns $cheevos";
}

is $dbh.execute(
    "SELECT earned FROM save_solution(3, 1, '⛳', 'π', 'c', 1)",
).row, '{different-strokes}', 'Earns {different-strokes}';

for $dbh.execute('SELECT unnest(enum_range(null::lang))').allrows.flat {
    my $earns = %langs{ my $i = ++$ } // '{}';

    # Add hole-specific cheevos on the end.
    $earns.=subst: '{', '{sounds-quite-nice,' if $_ eq 'd';
    $earns.=subst: '{', '{caffeinated,'       if $_ eq 'javascript';
    $earns.=subst: '{', '{just-kidding,'      if $_ eq 'k';
    $earns.=subst: '{', '{tim-toady,'         if $_ eq 'raku';
    $earns.=subst: ',}', '}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, ?, 'ab', 'musical-chords', ?, 1)",
        $_ eq 'assembly' ?? Nil !! 2, $_,
    ).row, $earns, "Lang $i ($_) earns $earns";
}

done-testing;
