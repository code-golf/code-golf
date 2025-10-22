use t;

constant %holes = <
      1 {hello-world}
      4 {fore}
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
     90 {right-on}
     99 {neunundneunzig-luftballons}
    100 {centenarian}
    107 {busy-beaver}
    111 {disappearing-act}
    120 {five}
>;

constant %langs = <
    12 {polyglot}
    24 {polyglutton}
    36 {omniglot}
    48 {omniglutton}
>;

my $dbh     = dbh;
my $session = new-golfer :$dbh;

$client.get: 'https://app/about',
    headers => { cookie => "__Host-session=$session" };

is $dbh.execute('SELECT ARRAY(SELECT cheevo FROM cheevos)').row, '{rtfm}',
    'GET /about earns {rtfm}';

for $dbh.execute('SELECT id FROM holes WHERE experiment = 0').allrows.flat {
    my $cheevos = %holes{ my $i = ++$ } // '{}';

    # Add hole-specific cheevos to the start.
    $cheevos.=subst: '{', '{interview-ready,' if $_ eq 'fizz-buzz';
    $cheevos.=subst: '{', '{solve-quine,'     if $_ eq 'quine';
    $cheevos.=subst: ',}', '}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', ?, 'c', 1::smallint, 1)", $_,
    ).row, $cheevos, "$_/c earns $cheevos";
}

for <
    24-game                          tex        {texnical-know-how}
    brainfuck                        brainfuck  {inception}
    divisors                         php        {elephpant-in-the-room}
    css-colors                       basic      {horse-of-a-different-color}
    evil-numbers                     scheme     {evil-scheme}
    factorial-factorisation          factor     {x-factor}
    game-of-life                     elixir     {alchemist}
    hexdump                          hexagony   {hextreme-agony}
    look-and-say                     sed        {simon-sed}
    pascals-triangle                 pascal     {under-pressure}
    poker                            fish       {fish-n-chips}
    quine                            python     {ouroboros}
    rijndael-s-box                   c-sharp    {s-box-360}
    rock-paper-scissors-spock-lizard janet      {dammit-janet}
    seven-segment                    assembly   {assembly-required}
    si-units                         powershell {watt-are-you-doing}
    star-wars-opening-crawl          tex        {typesetter}
    sudoku                           hexagony   {off-the-grid}
    ten-pin-bowling                  cobol      {cobowl}
    united-states                    pascal     {going-postal}
    𝑒                                r          {emergency-room}
> -> $hole, $lang, $cheevos {
    is $dbh.execute(
        "SELECT earned FROM save_solution(2, ?, 'ab', ?, ?, 1::smallint, 1)",
        $lang eq 'assembly' ?? Nil !! 2, $hole, $lang,
    ).row, $cheevos, "$hole/$lang earns $cheevos";
}

for <brainfuck d hexagony javascript nim swift sql zig> {
    my $cheevos = $_ eq 'zig' ?? '{pangramglot}' !! '{}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, 2, 'ab', 'pangram-grep', ?, 1::smallint, 1)",
        $_,
    ).row, $cheevos, "pangram-grep/$_ earns $cheevos";
}

is $dbh.execute(
    "SELECT earned FROM save_solution(3, 1, '⛳', 'π', 'c', 1::smallint, 1)",
).row, '{different-strokes}', 'Earns {different-strokes}';

for $dbh.execute('SELECT id FROM langs WHERE experiment = 0').allrows.flat {
    my $earns = %langs{ my $i = ++$ } // '{}';

    # Add hole-specific cheevos on the front.
    $earns.=subst: '{', '{sounds-quite-nice,' if $_ eq 'd';
    $earns.=subst: '{', '{piña-colada,'       if $_ eq 'elixir';
    $earns.=subst: '{', '{caffeinated,'       if $_ eq 'java';
    $earns.=subst: '{', '{go-forth,'          if $_ eq 'go';
    $earns.=subst: '{', '{just-kidding,'      if $_ eq 'k';
    $earns.=subst: '{', '{tim-toady,'         if $_ eq 'raku';
    $earns.=subst: ',}', '}';

    is $dbh.execute(
        "SELECT earned FROM save_solution(2, ?, 'ab', 'musical-chords', ?, 1::smallint, 1)",
        $_ eq 'assembly' ?? Nil !! 2, $_,
    ).row, $earns, "Lang $i ($_) earns $earns";
}

done-testing;
