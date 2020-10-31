use Cro::HTTP::Client;
use DBIish;
use Test;
use TOML::Thumb;

# Export Test & TOML::Thumb to caller.
sub EXPORT { %( Test::EXPORT::DEFAULT::, TOML::Thumb::EXPORT::DEFAULT:: ) }

# Block until the app is up or we time out.
react {
    whenever Promise.in(10) { bail-out 'Timed our waiting for app' }
    whenever start {
        sleep .1 until try await IO::Socket::Async.connect: 'app', 1080;
    } { done }
};

unit module t;

our $client is export = Cro::HTTP::Client.new: :ca({ :insecure }), :http(1.1);

sub dbh is export {
    my $dbh = DBIish.connect: 'Pg';

    $dbh.execute: 'SET client_min_messages TO WARNING';
    $dbh.execute: 'TRUNCATE code, users RESTART IDENTITY CASCADE';

    $dbh;
}

sub post-solution(:$code, :$hole = 'fizz-buzz', :$lang = 'raku', :$session = '') is export {
    my $res = await $client.post: 'https://app:1443/solution',
        content-type => 'application/json',
        headers      => [ cookie => "__Host-session=$session" ],
        body         => { Code => $code, Hole => $hole, Lang => $lang };

    await $res.body;
}
