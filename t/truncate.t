use t;

my $res = run 'say "a" x 128 and say STDERR "b" x 128 for 0..1024';
is $res<runs>[0]<stderr>.chars, 128 * 1024, 'Stderr is limited to 128 KiB';
is $res<runs>[0]<stdout>.chars, 128 * 1024, 'Stdout is limited to 128 KiB';

is run('say "abc", " " x 2**16  ')<runs>[0]<stdout>, 'abc', 'Really long line';
is run('say "abc", " \n\r\t" x 9')<runs>[0]<stdout>, 'abc', 'Lots of trailing whitespace';
is run('say "abc \t \x0B \f \r "')<runs>[0]<stdout>, 'abc', 'Different whitespace';

my $code = 'say "Fizz" x !($_ % 3) . "Buzz" x !($_ % 5) || $_ for 1 .. 100';

is-deeply run($code)<runs>[0]<pass>:kv.Hash,
    { :pass }, 'Passing without overflow';

is-deeply run("$code; say ' ' x (128 * 1024)")<runs>[0]<pass>:kv.Hash,
    { :!pass }, 'Failing with overflow';

sub run { post-solution :$^code, :lang<perl> }

done-testing;
