use t;

for 'config/langs.toml'.IO.&from-toml.map({
    .key.lc.subst(' ', '-').trans( qw[# + ><>] => qw[-sharp p fish] ) => .value<example>;
}).sort -> (:key($lang), :value($code)) {
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

            # Pascal & TeX prints lots of info to STDERR.
            is $got<Out>, $exp, 'Out';
            is $got<Err>,   '', 'Err' if $lang ne 'pascal' | 'tex';
        }
    }
}

done-testing;
