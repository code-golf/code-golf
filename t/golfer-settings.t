use t;

throws-like { $client.get: 'https://app/golfer/settings' },
    Exception, message => /'403 Forbidden'/, '403 with no session';

my $dbh     = dbh;
my $session = new-golfer :$dbh :id(123) :name<foo>;
              new-golfer :$dbh :id(456) :name<bar>;

is-deeply settings, {
    :country(Str), :layout<default>, :keymap<default>, :referrer_id(Int),
    :!show_country, :theme<auto>, :time_zone(Str),
}, 'DB has defaults';

for <country layout keymap time_zone> {
    throws-like { post %( :country(''), :layout<default>, :keymap<default>, :time_zone(''), $_ => 'baz' ) },
        Exception, message => /'400 Bad Request'/, "400 with invalid $_";
}

post my %args = :country<GB>, :layout<tabs>, :keymap<vim>, :show_country<on>,
    :theme<dark>, :time_zone<Europe/London>;

is-deeply settings, %( |%args, :referrer_id(Int), :show_country ),
    'DB is updated';

post %( |%args, :referrer<BaR> );   # Case-insensitive

is-deeply settings<referrer_id>, 456, 'referrer_id is 456';

$dbh.execute: 'DELETE FROM users WHERE id = 456';   # ON DELETE SET NULL

is-deeply settings<referrer_id>, Int, 'referrer_id is NULL';

post %( |%args, :referrer<foo> );   # Can't be yourself

is-deeply settings<referrer_id>, Int, 'referrer_id is NULL';

sub settings { $dbh.execute(q:to/SQL/).row :hash }
    SELECT country, layout, keymap, referrer_id, show_country, theme, time_zone
      FROM users
     WHERE id = 123
SQL

sub post(%content) {
    $client.post: 'https://app/golfer/settings', :%content,
        headers => { cookie => "__Host-session=$session" };
}

done-testing;
