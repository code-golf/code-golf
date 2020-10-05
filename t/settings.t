use t;

throws-like { await $client.get: 'https://app:1443/golfer/settings' },
    Exception, message => /'Server responded with 403 Forbidden'/,
    '403 with no session';

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'test')";

is $dbh.execute('SELECT time_zone FROM users').row.head, 'UTC',
    'Time zone defaults to "UTC"';

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

await $client.post: 'https://app:1443/golfer/settings',
    content-type => 'application/x-www-form-urlencoded',
    headers      => [ cookie => "__Host-session=$session" ],
    body         => {
        time_zone => 'Europe/London',
        keymap    => 'vim',
    };

is $dbh.execute('SELECT time_zone FROM users').row.head, 'Europe/London',
    'Time zone updated to "Europe/London"';

is $dbh.execute('SELECT keymap FROM users').row.head, 'vim',
    'Keymap preference updated to "vim"';

done-testing;
