use HTTP::Tiny;
use JSON::PP;
use Test2::V0;

my $res = HTTP::Tiny->new->post(
    'https://code.golf/solution',
    {   content => encode_json {
            Code => 'say "a" x 1024 and say STDERR "b" x 1024 for 0..128',
            Hole => 'fizz-buzz',
            Lang => 'perl',
        },
    },
);

die $res->{content} unless $res->{success};

$res = decode_json $res->{content};

is length $res->{Err}, 128 * 1024, 'Stderr is limited to 128 KiB';
is length $res->{Out}, 128 * 1024, 'Stdout is limited to 128 KiB';

done_testing;
