use t;

constant %trophies = <
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

for from-toml slurp 'holes.toml' {
    next if .value<experiment>;     # Experimental holes can't be saved.
    next if .key eq 'Fizz Buzz';    # This is tested lower.

    my $hole     = .key.lc.trans: ' ’' => '-', :d;
    my $trophies = %trophies{ my $i = ++$ } // '{}';

    is $dbh.execute("SELECT save_solution('', ?, 'c', 1)", $hole).row,
        $trophies, "Solution $i earns $trophies";
}

for <
    brainfuck       brainfuck {inception}
    divisors        php       {elephpant-in-the-room}
    fizz-buzz       haskell   {bullseye,interview-ready}
    quine           python    {ouroboros}
    ten-pin-bowling cobol     {cobowl}
> -> $hole, $lang, $trophies {
    is $dbh.execute("SELECT save_solution('', ?, ?, 1)", $hole, $lang).row,
        $trophies, "$hole/$lang earns $trophies";
}

done-testing;
