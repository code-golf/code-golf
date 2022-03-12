use t;
use HTTP::Tiny;

dbh.execute: "INSERT INTO users (id,  login)
                         VALUES (123, 'foo'), (456, 'bar'), (789, 'baz')";

my $ua = HTTP::Tiny.new :!max-redirect;

for «
    200 /api                          ｢ 'openapi:' .*?                         ｣
    200 /api/cheevos                  ｢ '[' .*? '"RTFM"' .*? '"tl;dr"' .*? ']' ｣
    200 /api/cheevos/tl-dr            ｢ '{' .*?              '"tl;dr"' .*? '}' ｣
    404 /api/cheevos/unknown          ｢ 'null'                                 ｣
    200 /api/langs                    ｢ '[' .*? '"Perl"' .*? '"Raku"'  .*? ']' ｣
    200 /api/langs/raku               ｢ '{' .*?              '"Raku"'  .*? '}' ｣
    404 /api/langs/unknown            ｢ 'null'                                 ｣
    404 /api/not-found                ｢ 'null'                                 ｣
    500 /api/panic                    ｢ 'null'                                 ｣
    200 /api/suggestions/golfers      ｢ '["bar", "baz", "foo"]'                ｣
    200 /api/suggestions/golfers?q=ba ｢ '["bar", "baz"]'                       ｣
    200 /api/suggestions/golfers?q=x  ｢ '[]'                                   ｣
    200 /api/suggestions/golfers?q=z  ｢ '["baz"]'                              ｣
» -> $status, $path, $content {
    subtest $path, {
        my %res = $ua.get: "https://app:1443$path";

        is %res<status>, $status, 'status';
        is %res<headers><access-control-allow-origin>, '*', 'CORS';

        # /api is YML, not JSON.
        my $content-type = $path eq '/api' ?? 'text/plain; charset=utf-8'
                                           !! 'application/json';
        is %res<headers><content-type>, $content-type, 'content-type';

        like %res<content>.decode, /^ <$content> \n? $/, 'content';
    }
}

done-testing;
