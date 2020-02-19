use feature 'state';

use HTTP::Tiny;
use JSON::PP;
use Test2::V0;

is run('say %*ENV'), '{}', 'Environment';

is run('say $*KERNEL.hostname'), 'code-golf', 'Hostname';

is [ map [ split / / ], split /\n/, run('say slurp "/proc/mounts"') ] => [
    [ 'overlay', '/', 'overlay', match('^ro,'), 0, 0 ],
    [ 'none', '/proc', 'proc', 'ro,relatime', 0, 0 ],
] => '/proc/mounts';

like { run('say slurp "/proc/self/status"') =~ /(.+):\s*(.*)/g } => {
    # No capabilities
    CapAmb => '0000000000000000',
    CapBnd => '0000000000000000',
    CapEff => '0000000000000000',
    CapInh => '0000000000000000',
    CapPrm => '0000000000000000',

    # User/Group nobody
    Gid => "65534\t65534\t65534\t65534",
    Uid => "65534\t65534\t65534\t65534",

    Name => 'raku',

    NoNewPrivs => 1,

    Pid => 1,
    PPid => 0,

    # Seccomp filter mode
    Seccomp => 2,

    Speculation_Store_Bypass => 'thread force mitigated',
} => '/proc/self/status';

sub run {
    my $res = ( state $ua = HTTP::Tiny->new )->post(
        'https://code-golf.io/solution',
        { content => encode_json { Code => $_[0], Hole => 'π', Lang => 'raku' } },
    );

    die $res->{content} unless $res->{success};

    decode_json( $res->{content} )->{Out};
}

done_testing;
