use HTTP::Tiny;
use JSON::PP;
use Test2::V0;

my $res = HTTP::Tiny->new->post(
    'https://code.golf/solution',
    {   content => encode_json {
            Code => 'sleep 8',
            Hole => 'fizz-buzz',
            Lang => 'perl',
        },
    },
);

die $res->{content} unless $res->{success};

$res = decode_json $res->{content};

is $res->{Err}, 'Killed for exceeding the 7s timeout.';

is int $res->{Took} / 1e9, 7;

done_testing;
