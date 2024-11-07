use DBIish;
use HTTP::Tiny;
use JSON::Fast;
use Test;
use TOML::Thumb;

# Export Test & TOML::Thumb to caller.
sub EXPORT { %( Test::EXPORT::DEFAULT::, TOML::Thumb::EXPORT::DEFAULT:: ) }

unit module t;

our $client is export = HTTP::Tiny.new :throw-exceptions;

# Connect to the DB, truncate solutions & users, return the handle.
sub dbh is export {
    my $dbh = DBIish.connect: 'Pg';

    $dbh.execute: 'SET client_min_messages TO WARNING';
    $dbh.execute: 'TRUNCATE solutions, users RESTART IDENTITY CASCADE';

    $dbh;
}

# Create a new Golfer, log them in, and return the session ID.
# If called with no DB handle then it'll implicitly connect & truncate tables.
sub new-golfer(:$dbh = dbh, :$id = 1, :$name = 'Bob') is export {
    $dbh.execute('INSERT INTO users (id, login) VALUES ($1, $2)', $id, $name);
    $dbh.execute('INSERT INTO sessions (user_id) VALUES ($1) RETURNING id', $id).row.head;
}

# Submit a solution and return the deserialised response.
sub post-solution(:$code, :$hole = 'fizz-buzz', :$lang = 'raku', :$session = '') is export {
    $client.post(
        'https://app/solution',
        content => to-json({ Code => $code, Hole => $hole, Lang => $lang }),
        headers => { cookie => "__Host-session=$session" },
    )<content>.decode.&from-json;
}
