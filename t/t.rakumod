use DBIish;
use HTTP::Tiny;
use JSON::Fast;
use Test;
use TOML::Thumb;

# Export Test & TOML::Thumb to caller.
sub EXPORT { %( Test::EXPORT::DEFAULT::, TOML::Thumb::EXPORT::DEFAULT:: ) }

unit module t;

# TODO Remove :!max-redirect once https://gitlab.com/jjatria/http-tiny/-/issues/12
our $client is export = HTTP::Tiny.new :!max-redirect :throw-exceptions;

sub createSession($dbh, int $userId) is export {
    $dbh.execute("INSERT INTO sessions (user_id) VALUES ($userId) RETURNING id").row.head;
}

sub createUser($dbh, int $id) is export {
    $dbh.execute: "INSERT INTO users (id, login) VALUES ($id, 'Bob')";
}

sub dbh is export {
    my $dbh = DBIish.connect: 'Pg';

    $dbh.execute: 'SET client_min_messages TO WARNING';
    $dbh.execute: 'TRUNCATE solutions, users RESTART IDENTITY CASCADE';

    $dbh;
}

sub post-solution(:$code, :$hole = 'fizz-buzz', :$lang = 'raku', :$session = '') is export {
    $client.post(
        'https://app/solution',
        content => to-json({ Code => $code, Hole => $hole, Lang => $lang }),
        headers => { cookie => "__Host-session=$session" },
    )<content>.decode.&from-json;
}
