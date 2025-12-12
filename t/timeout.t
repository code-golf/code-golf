use t;

# Note timeout is 5s in live, 10s under e2e. See hole/play.go for details.

subtest 'Timeout' => {
    my $res = post-solution( :code('sleep 11') )<runs>[0];

    is $res<stderr>, 'Killed for exceeding the 10s timeout.', 'Correct error';

    is floor( $res<time_ns> / 1e9 ), 10, 'Correct time';
};

subtest 'Timeout with correct output' => {
    my $res = post-solution( :code(q:to/CODE/) )<runs>[0];
        say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1 .. 100;
        $*OUT.flush;
        sleep 11;
    CODE

    is $res<stderr>, 'Killed for exceeding the 10s timeout.', 'Correct error';

    is $res<stdout>, slurp('config/hole-answers/fizz-buzz.txt').trim, 'Correct output';

    nok $res<pass>, 'Solution fails';
};

done-testing;
