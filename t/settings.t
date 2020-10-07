use t;

throws-like { await $client.get: 'https://app:1443/golfer/settings' },
    Exception, message => /'Server responded with 403 Forbidden'/,
    '403 with no session';

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'test')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

my %defaults = :keymap<default>, :time_zone<UTC>;

is settings, %defaults, 'DB has defaults';

for <country keymap time_zone> {
    throws-like { post %( |%defaults, $_ => 'foo' ) },
        Exception, message => /'Server responded with 400 Bad Request'/,
        "400 with invalid $_";
}

post my %args = :keymap<vim>, :time_zone<Europe/London>;

is settings, %args, 'DB is updated';

sub settings {
    $dbh.execute('SELECT keymap, time_zone FROM users').row(:hash);
}

sub post(%body) {
    await $client.post: 'https://app:1443/golfer/settings',
        body         => %body,
        content-type => 'application/x-www-form-urlencoded',
        headers      => [ cookie => "__Host-session=$session" ];
}

done-testing;
