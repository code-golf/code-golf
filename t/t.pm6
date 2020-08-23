use Cro::HTTP::Client;
use Test;

# Export Test to caller.
sub EXPORT { Test::EXPORT::DEFAULT:: }

# Block until the app is up or we time out.
react {
    whenever Promise.in(3) { bail-out 'Timed our waiting for app' }
    whenever start {
        sleep .1 until try await IO::Socket::Async.connect: 'app', 1080;
    } { done }
};

unit module t;

sub post-solution(:$code, :$hole = 'fizz-buzz', :$lang = 'raku') is export {
    state $client = Cro::HTTP::Client.new: ca => { :insecure };

    my $res = await $client.post: 'https://app:1443/solution',
        content-type => 'application/json',
        body         => { Code => $code, Hole => $hole, Lang => $lang };

    await $res.body;
}
