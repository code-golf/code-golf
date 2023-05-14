use t;

my $res = run 'say "a" x 1024 and say STDERR "b" x 1024 for 0..128';
is $res<Err>.chars, 128 * 1024, 'Stderr is limited to 128 KiB';
is $res<Out>.chars, 128 * 1024, 'Stdout is limited to 128 KiB';

is run('say "abc", " " x 2**16  ')<Out>, 'abc', 'Really long line';
is run('say "abc", " \n\r\t" x 9')<Out>, 'abc', 'Lots of trailing whitespace';
is run('say "abc \t \x0B \f \r "')<Out>, 'abc', 'Different whitespace';

sub run { post-solution :$^code, :lang<perl> }

done-testing;
