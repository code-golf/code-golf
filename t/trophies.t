use t;
use TOML::Thumb;

constant %trophies = <
    1  {hello-world}
    13 {bakers-dozen}
    19 {the-watering-hole}
    42 {dont-panic}
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
    brainfuck brainfuck {inception}
    divisors  php       {elephpant-in-the-room}
    fizz-buzz haskell   {interview-ready}
    quine     python    {ouroboros}
> -> $hole, $lang, $trophies {
    is $dbh.execute("SELECT save_solution('', ?, ?, 1)", $hole, $lang).row,
        $trophies, "$hole/$lang earns $trophies";
}

done-testing;
