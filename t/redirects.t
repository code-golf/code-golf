use HTTP::Tiny;
use Test;

my constant $ua = HTTP::Tiny.new :!max-redirect;
my constant $query = '?foo=bar';

# Redirects.
for <
    GET    /api/holes/billiard                    /api/holes/billiards
    GET    /api/langs/perl6                       /api/langs/raku
    DELETE /api/notes/billiard/perl6              /api/notes/billiards/raku
    GET    /api/notes/billiard/perl6              /api/notes/billiards/raku
    PUT    /api/notes/billiard/perl6              /api/notes/billiards/raku
    GET    /billiard                              /billiards
    GET    /rankings/recent-holes/perl6/bytes     /rankings/recent-holes/raku/bytes
    GET    /recent/perl6                          /recent/raku
    GET    /recent/solutions/billiard/perl6/bytes /recent/solutions/billiards/raku/bytes
    GET    /scores/billiard/perl6                 /scores/billiards/raku
> -> $method, $start, $end {
    my $res = $ua.request: $method, "https://app:443$start$query";

    subtest "$method $start → $end" => sub {
        is $res<headers><location>, $end ~ $query, 'Location header';
        is $res<status>, 308, 'Status';
    }
}

# Aliases.
for <
    GET /lambda /%ce%bb
    GET /pi     /%cf%80
    GET /tau    /%cf%84
> -> $method, $start, $end {
    my $res = $ua.request: $method, "https://app:443$start$query";

    subtest "$method $start → $end" => sub {
        is $res<headers><location>, $end ~ $query, 'Location header';
        is $res<status>, 307, 'Status';
    }
}

done-testing;
