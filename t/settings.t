use t;

throws-like { await $client.get: 'https://app:1443/golfer/settings' },
    Exception, message => /'Server responded with 403 Forbidden'/,
    '403 with no session';

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'test')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

is-deeply settings,
    { :country(Str), :keymap<default>, :!show_country, :time_zone(Str) },
    'DB has defaults';

for <country keymap time_zone> {
    throws-like { post %( :country(''), :keymap<default>, :time_zone(''), $_ => 'foo' ) },
        Exception, message => /'Server responded with 400 Bad Request'/,
        "400 with invalid $_";
}

post my %args = :country<GB>, :keymap<vim>, :show_country<on>, :time_zone<Europe/London>;

is-deeply settings, %( |%args, :show_country ), 'DB is updated';

sub settings {
    $dbh.execute('SELECT country, keymap, show_country, time_zone FROM users').row(:hash);
}

sub post(%body) {
    await $client.post: 'https://app:1443/golfer/settings',
        body         => %body,
        content-type => 'application/x-www-form-urlencoded',
        headers      => [ cookie => "__Host-session=$session" ];
}

done-testing;
