use Cro::HTTP::Client;
use DBIish;
use Test;
use TOML::Thumb;

# Export Test & TOML::Thumb to caller.
sub EXPORT { %( Test::EXPORT::DEFAULT::, TOML::Thumb::EXPORT::DEFAULT:: ) }

unit module t;

our $client is export = Cro::HTTP::Client.new: :ca({ :insecure }), :http(1.1);

sub dbh is export {
    my $dbh = DBIish.connect: 'Pg';

    $dbh.execute: 'SET client_min_messages TO WARNING';
    $dbh.execute: 'TRUNCATE solutions, users RESTART IDENTITY CASCADE';

    $dbh;
}

sub post-solution(:$code, :$hole = 'fizz-buzz', :$lang = 'raku', :$session = '') is export {
    my $res = await $client.post: 'https://app:1443/solution',
        content-type => 'application/json',
        headers      => [ cookie => "__Host-session=$session" ],
        body         => { Code => $code, Hole => $hole, Lang => $lang };

    await $res.body;
}
