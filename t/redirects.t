use HTTP::Tiny;
use Test;

my constant $query = '?foo=bar';

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
    state $ua = HTTP::Tiny.new :!max-redirect;

    my $location = $ua.request(
        $method, "https://app:443$start$query",
    )<headers><location>;

    is $location, $end ~ $query, "$method $start → $end";
}

done-testing;
