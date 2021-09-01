use t;

for 'langs.toml'.IO.&from-toml.map({
    .key.lc.trans( qw[# ><>] => qw[-sharp fish] ) => .value<example>;
}).sort -> (:key($lang), :value($code)) {
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
            my $exp = join "\n", 'Hello, World!', |^10, |$got<Argv>;

            $exp ~= "\n" if $hole eq 'quine';

            # Pascal prints lots of info to STDERR.
            is $got<Out>, $exp, 'Out';
            is $got<Err>,   '', 'Err' if $lang ne 'pascal';
        }
    }
}

done-testing;
