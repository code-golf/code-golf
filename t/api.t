use t;
use HTTP::Tiny;

my $dbh = dbh;
$dbh.execute: "INSERT INTO users (id,  login)
                          VALUES (123, 'foo'), (456, 'bar'), (789, 'baz')";

$dbh.execute: "INSERT INTO notes VAlUES (123, 'π', 'c', '🥧')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

my $ua = HTTP::Tiny.new :!max-redirect;

for «
    200 /api                             ｢ 'openapi:' .*?                          ｣
    200 /api/cheevos                     ｢ '[' .*? '"RTFM"' .*? '"tl;dr"' .*? ']'  ｣
    200 /api/cheevos/tl-dr               ｢ '{' .*?              '"tl;dr"' .*? '}'  ｣
    404 /api/cheevos/unknown             ｢ 'null'                                  ｣
    200 /api/holes                       ｢ '[' .*? '"ISBN"' .*? '"π"'     .*? ']'  ｣
    200 /api/holes/π                     ｢ '{' .*?              '"π"'     .*? '}'  ｣
    404 /api/holes/unknown               ｢ 'null'                                  ｣
    200 /api/langs                       ｢ '[' .*? '"Perl"' .*? '"Raku"'  .*? ']'  ｣
    200 /api/langs/raku                  ｢ '{' .*?              '"Raku"'  .*? '}'  ｣
    404 /api/langs/unknown               ｢ 'null'                                  ｣
    404 /api/not-found                   ｢ 'null'                                  ｣
    200 /api/notes                       ｢ '[{"hole":"π","lang":"c","note":"🥧"}]' ｣
    200 /api/notes/π/c                   ｢ '🥧'                                    ｣
    404 /api/notes/π/d                   ｢ 'null'                                  ｣
    500 /api/panic                       ｢ 'null'                                  ｣
    200 /api/solutions-log?hole=π&lang=c ｢ '[]'                                    ｣
    200 /api/solutions-search?pattern=a  ｢ '[]'                                    ｣
    200 /api/suggestions/golfers         ｢ '["bar", "baz", "foo"]'                 ｣
    200 /api/suggestions/golfers?q=ba    ｢ '["bar", "baz"]'                        ｣
    200 /api/suggestions/golfers?q=x     ｢ '[]'                                    ｣
    200 /api/suggestions/golfers?q=z     ｢ '["baz"]'                               ｣
» -> $status, $path, $content {
    subtest $path, {
        my %res = $ua.get: "https://app$path",
            :headers( :cookie("__Host-session=$session") );

        is %res<status>, $status, 'status';
        is %res<headers><access-control-allow-origin>, '*', 'CORS';

        # Some routes return plain text, not JSON.
        my $content-type = $path eq '/api' | '/api/notes/π/c'
            ?? 'text/plain; charset=utf-8'
            !! 'application/json';
        is %res<headers><content-type>, $content-type, 'content-type';

        like %res<content>.decode, /^ <$content> \n? $/, 'content';
    }
}

done-testing;
