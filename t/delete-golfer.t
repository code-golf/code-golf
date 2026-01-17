# Test "ON DELETE" triggers when deleting a golfer.
# TODO Check it works without DELETE perms like solutions_log, weekly_solves.
use t;

my $dbh = dbh;

#########
# Setup #
#########

new-golfer :$dbh :id(1) :name<Alice>;
new-golfer :$dbh :id(2) :name<Bob>;

# Alice authored Fizz Buzz.
$dbh.execute: "INSERT INTO authors VALUES ('fizz-buzz', 1)";

# Alice referred Bob.
$dbh.execute: "UPDATE users SET referrer_id = 1 WHERE name = 'Bob'";

##########
# Before #
##########

is-deeply authors(), ({ hole => 'fizz-buzz', user_id => 1 },),
    'authors before delete';

is-deeply sessions(), ({ user_id => 1 }, { user_id => 2 }),
    'sessions before delete';

is-deeply users(), (
    { id => 1, name => 'Alice', referrer_id => Int },
    { id => 2, name => 'Bob',   referrer_id => 1   },
), 'users before delete';

##############
# Kill Alice #
##############

$dbh.execute: "DELETE FROM users WHERE name = 'Alice'";

#########
# After #
#########

is-deeply authors(), $(),
    'authors after delete';

is-deeply sessions(), ({ user_id => 2 },),
    'sessions after delete';

is-deeply users(), ({ id => 2, name => 'Bob', referrer_id => Int },),
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
    SELECT id, name, referrer_id FROM users
SQL
