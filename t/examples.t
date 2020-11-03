use t;

for sort from-toml slurp 'langs.toml' {
    my $lang = .key eq '><>' ?? 'fish' !! .key.lc.subst: '#', '-sharp';

    # TODO Remove this to ensure we have examples for every lang.
    # skip 'No example yet' unless .value<example>;
    next unless .value<example>;

    # Pick a hole that will definitely have unicode.
    my $res = post-solution
        code => .value<example>,
        hole => 'rock-paper-scissors-spock-lizard',
        lang => $lang;

    is $res<Out>, join("\n", 'Hello, World!', |^10, |$res<Argv>), $lang
        or diag $res<Err>;
}

done-testing;
