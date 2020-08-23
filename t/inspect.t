use t;

sub run($code) { post-solution(:$code)<Out> }

is run('say %*ENV'), '{}', 'Environment';

is run('say $*KERNEL.hostname'), 'code-golf', 'Hostname';

# TODO
# diag run('say slurp "/proc/mounts"').lines;
# [ 'overlay', '/', 'overlay', match('^ro,'), 0, 0 ]
# [ 'none', '/proc', 'proc', 'ro,relatime', 0, 0 ]

my %status   = run('say slurp "/proc/self/status"').split: / ':' \s+ | \n /;
my %expected = (
    # TODO No capabilities
    # CapAmb => '0000000000000000',
    # CapBnd => '0000000000000000',
    # CapEff => '0000000000000000',
    # CapInh => '0000000000000000',
    # CapPrm => '0000000000000000',

    # User/Group nobody
    Gid => "65534\t65534\t65534\t65534",
    Uid => "65534\t65534\t65534\t65534",

    Name => 'raku',

    NoNewPrivs => '1',

    Pid  => '1',
    PPid => '0',

    # Seccomp filter mode
    Seccomp => '2',
);

is-deeply %status{ %expected.keys }:kv.Hash, %expected, '/proc/self/status';

done-testing;
