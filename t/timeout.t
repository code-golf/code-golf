use t;

# Note timeout in 5s in live, 10s in E2E. See hole/play.go for details.

subtest 'Timeout' => {
    my $res = post-solution :code('sleep 11');

    is $res<Err>, 'Killed for exceeding the 10s timeout.', 'Correct error';

    is floor( $res<Took> / 1e9 ), 10, 'Correct took';
};

subtest 'Timeout with correct output' => {
    my $res = post-solution :code('say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1 .. 100; $*OUT.flush; sleep 11;');

    is $res<Err>, 'Killed for exceeding the 10s timeout.', 'Correct error';

    is $res<Out>, slurp('hole/answers/fizz-buzz.txt').trim, 'Correct output';

    nok $res<Pass>, 'Solution fails';
};

done-testing;
