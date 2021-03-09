use t;

for slurp("langs.toml").&from-toml.map({
    .key.lc.trans( qw[# ><>] => qw[-sharp fish] ) => .value<example>;
}).sort -> (:key($lang), :value($code)) {
    # TODO Remove this to ensure we have examples for every lang.
    next unless $code;

    # <built-in>: internal compiler error: Illegal instruction
    todo 'intermittent error' if $lang eq 'fortran';

    for (
        # Pick a hole that will definitely have unicode arguments.
        'rock-paper-scissors-spock-lizard',

        # Ensure PowerShell example works on Quine with explicit output.
        ( 'quine' if $lang eq 'powershell' ),
    ) -> $hole {
        subtest "$lang ($hole)" => {
            my $got = post-solution :$code :$hole :$lang;
            my $exp = join "\n", 'Hello, World!', |^10, |($got<Argv> // '');

            is $got<Out>, $exp, 'Out';
            is $got<Err>,   '', 'Err';
        }
    }
}

done-testing;
