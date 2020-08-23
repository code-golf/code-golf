use t;

my $res = post-solution
    code => 'say "a" x 1024 and say STDERR "b" x 1024 for 0..128',
    lang => 'perl';

is $res<Err>.chars, 128 * 1024, 'Stderr is limited to 128 KiB';
is $res<Out>.chars, 128 * 1024, 'Stdout is limited to 128 KiB';

done-testing;
