use t;
use HTTP::Tiny;

my $dbh = dbh;

$dbh.execute: "INSERT INTO users (id, login) VALUES (123, 'foo'), (456, 'bar')";

my $session = $dbh.execute(
    'INSERT INTO sessions (user_id) VALUES (123) RETURNING id').row.head;

is post('follow', 'foo')<status>, 400, "Can't follow yourself";
is post('unfollow', 'foo')<status>, 400, "Can't unfollow yourself";

is-deeply follows, {}, 'DB is empty';

is post('follow', 'bar')<status>, 303, 'Follow works';
is post('follow', 'bar')<status>, 303, 'Second follow does nothing';

is-deeply follows, { :follower_id(123), :followee_id(456) }, 'DB is updated';

is post('unfollow', 'bar')<status>, 303, 'Unfollow works';
is post('unfollow', 'bar')<status>, 303, 'Second unfollow does nothing';

is-deeply follows, {}, 'DB is empty';

sub follows { $dbh.execute('SELECT * FROM follows').row :hash }

sub post($action, $golfer) {
    state $ua = HTTP::Tiny.new :!max-redirect;

    $ua.post: "https://app:1443/golfers/$golfer/$action",
        headers => { cookie => "__Host-session=$session" };
}

done-testing;
