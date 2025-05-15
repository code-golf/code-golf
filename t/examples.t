use t;

for 'config/data/langs.toml'.IO.&from-toml.map({
    .key.lc.subst(' ', '-').trans( qw[# + ><>] => qw[-sharp p fish] ) => .value<example>;
}).sort -> (:key($lang), :value($code)) {
    for (
        # Pick a hole that will definitely have unicode arguments.
        'rock-paper-scissors-spock-lizard',

        # Ensure PowerShell example works on Quine with explicit output.
        ( 'quine' if $lang eq 'powershell' ),
    ) -> $hole {
        subtest "$lang ($hole)" => {
            my $got = post-solution( :$code :$hole :$lang )<runs>[0];
            my $exp = join "\n", 'Hello, World!', |^10, |$got<args>;

            $exp ~= "\n" if $hole eq 'quine';

            # Some langs print noise to STDERR, so don't test for no output.
            is $got<stdout>, $exp, 'Stdout';
            is $got<stderr>,   '', 'Stderr'
                if $lang ne any < factor kotlin tex vyxal >;
        }
    }
}

done-testing;
