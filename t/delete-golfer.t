# Test "ON DELETE" triggers when deleting a golfer.
use t;

my $dbh = dbh;

#########
# Setup #
#########

createUser($dbh, 1, 'Alice');
createUser($dbh, 2, 'Bob');

createSession($dbh, 1);
createSession($dbh, 2);

# Alice authored Fizz Buzz.
$dbh.execute: "INSERT INTO authors VALUES ('fizz-buzz', 1)";

# Alice referred Bob.
$dbh.execute: "UPDATE users SET referrer_id = 1 WHERE login = 'Bob'";

##########
# Before #
##########

is-deeply authors(), ({ hole => 'fizz-buzz', user_id => 1 },),
    'authors before delete';

is-deeply sessions(), ({ user_id => 1 }, { user_id => 2 }),
    'sessions before delete';

is-deeply users(), (
    { id => 1, login => 'Alice', referrer_id => Int },
    { id => 2, login => 'Bob',   referrer_id => 1   },
), 'users before delete';

##############
# Kill Alice #
##############

$dbh.execute: "DELETE FROM users WHERE login = 'Alice'";

#########
# After #
#########

is-deeply authors(), $(),
    'authors after delete';

is-deeply sessions(), ({ user_id => 2 },),
    'sessions after delete';

is-deeply users(), ({ id => 2, login => 'Bob', referrer_id => Int },),
    'users after delete';

done-testing;

#########
# Utils #
#########

sub authors { $dbh.execute(q:to/SQL/).allrows :array-of-hash }
    SELECT * FROM authors
SQL

sub sessions { $dbh.execute(q:to/SQL/).allrows :array-of-hash }
    SELECT user_id FROM sessions
SQL

sub users { $dbh.execute(q:to/SQL/).allrows :array-of-hash }
    SELECT id, login, referrer_id FROM users
SQL
