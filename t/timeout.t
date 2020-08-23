use t;

my $res = post-solution :code('sleep 8');

is $res<Err>, 'Killed for exceeding the 7s timeout.', 'Correct error';

is floor( $res<Took> / 1e9 ), 7, 'Correct took';

done-testing;
