use t;
use HTTP::Tiny;

new-golfer;

my constant $ua = HTTP::Tiny.new :!max-redirect;
my constant $query = '?foo=bar';

# Redirects.
for <
    GET    /api/holes/billiard                   /api/holes/billiards
    GET    /api/langs/perl6                      /api/langs/raku
    DELETE /api/notes/billiard/perl              /api/notes/billiards/perl
    GET    /api/notes/billiards/perl6            /api/notes/billiards/raku
    PUT    /api/notes/billiard/go                /api/notes/billiards/go
    GET    /billiard                             /billiards
    GET    /golfers/Bob/isbn/perl6/bytes         /golfers/Bob/isbn/raku/bytes
    GET    /rankings/cheevos                     /rankings/cheevos/all
    GET    /rankings/recent-holes/perl6/bytes    /rankings/recent-holes/raku/bytes
    GET    /recent                               /recent/solutions/all/all/bytes
    GET    /recent/perl6                         /recent/raku
    GET    /recent/solutions/billiard/perl/bytes /recent/solutions/billiards/perl/bytes
    GET    /scores/billiard/perl                 /scores/billiards/perl
    GET    /users/Bob                            /golfers/Bob
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
