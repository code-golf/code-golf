use t;

throws-like { $client.get: 'https://app/golfer/settings' },
    Exception, message => /'403 Forbidden'/, '403 with no session';

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'foo'), (456, 'bar')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

is-deeply settings, {
    :country(Str), :referrer_id(Int), :!show_country,
    :theme<auto>,  :time_zone(Str),
}, 'DB has defaults';

for <country time_zone> {
    throws-like { post %( :country(''), :time_zone(''), $_ => 'baz' ) },
        Exception, message => /'400 Bad Request'/, "400 with invalid $_";
}

post my %args = :country<GB>, :show_country<on>, :theme<dark>,
    :time_zone<Europe/London>;

is-deeply settings, %( |%args, :referrer_id(Int), :show_country ),
    'DB is updated';

post %( |%args, :referrer<BaR> );   # Case-insensitive

is-deeply settings<referrer_id>, 456, 'referrer_id is 456';

$dbh.execute: 'DELETE FROM users WHERE id = 456';   # ON DELETE SET NULL

is-deeply settings<referrer_id>, Int, 'referrer_id is NULL';

post %( |%args, :referrer<foo> );   # Can't be yourself

is-deeply settings<referrer_id>, Int, 'referrer_id is NULL';

sub settings { $dbh.execute(q:to/SQL/).row :hash }
    SELECT country, referrer_id, show_country, theme, time_zone
      FROM users
     WHERE id = 123
SQL

sub post(%content) {
    $client.post: 'https://app/golfer/settings', :%content,
        headers => { cookie => "__Host-session=$session" };
}

done-testing;
