use t;

sub run($code) { post-solution(:$code)<Out> }

is run('say %*ENV'), '{}', 'Environment';

is run('say $*KERNEL.hostname'), 'code-golf', 'Hostname';

my $mounts = run('say slurp "/proc/mounts"');

$mounts ~~ s:g/ ',inode64'       //;  # Older kernels wont have this flag.
$mounts ~~ s:g/ ',lowerdir=' \S+ //;  # Strip the Docker specific crap.

# Note on live root is read-only but doing that on dev breaks go build.
is $mounts, Q:to/MOUNTS/.chomp, '/proc/mounts';
    overlay / overlay rw,relatime 0 0
    tmpfs /dev tmpfs rw,nosuid,relatime 0 0
    proc /proc proc ro,nosuid,nodev,noexec,relatime 0 0
    tmpfs /tmp tmpfs rw,nosuid,nodev,relatime 0 0
    MOUNTS

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
