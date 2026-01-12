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
    $cheevos.=subst: '{', '{smörgåsbord,'     if $_ eq 'catalans-constant';
    $cheevos.=subst: '{', '{interview-ready,' if $_ eq 'fizz-buzz';
    $cheevos.=subst: '{', '{solve-quine,'     if $_ eq 'quine';
    $cheevos.=subst: ',}', '}';

    is save-solution(:hole($_) :lang<c>), $cheevos, "$_/c earns $cheevos";
}

for Q:ww<
    {alchemist} game-of-life elixir
    {alphabet-soup} scrambled-sort 'c d j'
    {archivist} isbn 'basic cobol common-lisp'
    {assembly-required} seven-segment assembly
    {bird-is-the-word} levenshtein-distance 'awk prolog sql'
    {cobowl} ten-pin-bowling cobol
    {count-to-ten} collatz 'c jq fish perl janet pascal haskell hexagony brainfuck javascript'
    {dammit-janet} rock-paper-scissors-spock-lizard janet
    {elephpant-in-the-room} divisors php
    {emergency-room} 𝑒 r
    {evil-scheme} evil-numbers scheme
    {fish-n-chips} poker fish
    {flag-those-mines} minesweeper 'f-sharp factor forth'
    {full-stack-dev} css-colors 'javascript php sql'
    {going-postal} united-states pascal
    {happy-go-lucky} 'happy-numbers lucky-numbers' go
    {hextreme-agony} hexdump hexagony
    {horse-of-a-different-color} css-colors basic
    {inception} brainfuck brainfuck
    {jeweler} diamonds 'crystal ruby'
    {mary-had-a-little-lambda} λ 'clojure coconut common-lisp'
    {never-eat-shredded-wheat} arrows 'nim elixir sed wren'
    {off-the-grid} sudoku hexagony
    {ouroboros} quine python
    {pangramglot} pangram-grep 'brainfuck d hexagony javascript nim swift sql zig'
    {ring-toss} tower-of-hanoi 'cobol factor fortran go groovy kotlin ocaml prolog python'
    {s-box-360} rijndael-s-box c-sharp
    {simon-sed} look-and-say sed
    {sinosphere} mahjong 'c j v'
    {sounds-quite-nice} musical-chords 'c c-sharp d'
    {texnical-know-how} 24-game tex
    {typesetter} star-wars-opening-crawl tex
    {under-pressure} pascals-triangle pascal
    {watt-are-you-doing} si-units powershell
    {when-in-rome} roman-to-arabic 'c d v'
    {x-factor} factorial-factorisation factor
    {zoodiac-signs} zodiac-signs 'awk basic civet'
> -> $cheevos, $holes, $langs {
    my @holes = $holes.words;
    my @langs = $langs.words;

    for @holes -> $hole {
        for @langs -> $lang {
            my $earns = $hole eq @holes.tail && $lang eq @langs.tail
                ?? $cheevos !! '{}';

            is save-solution(:$hole :$lang), $earns, "$hole/$lang earns $earns";
        }
    }
}

is save-solution(:code<⛳> :hole<π> :lang<c>), '{different-strokes}',
    'Earns {different-strokes}';

for $dbh.execute('SELECT id FROM langs WHERE experiment = 0').allrows.flat {
    my $earns = %langs{ my $i = ++$ } // '{}';

    # Add hole-specific cheevos on the front.
    $earns.=subst: '{', '{piña-colada,'       if $_ eq 'elixir';
    $earns.=subst: '{', '{caffeinated,'       if $_ eq 'java';
    $earns.=subst: '{', '{go-forth,'          if $_ eq 'go';
    $earns.=subst: '{', '{just-kidding,'      if $_ eq 'k';
    $earns.=subst: '{', '{tim-toady,'         if $_ eq 'raku';
    $earns.=subst: '{', '{down-to-the-metal,' if $_ eq 'rust';
    $earns.=subst: ',}', '}';

    is save-solution(:hole<cubes> :lang($_)), $earns,
        "Lang $i ($_) earns $earns";
}

done-testing;

sub save-solution(:$code = 'ab', :$hole, :$lang) {
    return $dbh.execute(
        'SELECT earned FROM save_solution(
            bytes := ?, chars := ?, code := ?, hole := ?,
            lang  := ?, time_ms := 1::smallint, user_id := 1)',
        $code.encode('UTF-8').bytes, $lang eq 'assembly' ?? Nil !! $code.chars,
        $code, $hole, $lang,
    ).row;
}
