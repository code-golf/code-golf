use HTTP::Tiny;
use Test2::V0;

use constant PATH => '/foo?bar=baz';

my $ua = HTTP::Tiny->new;

is $ua->get( $_ . PATH )->{url}, 'https://code.golf' . PATH, $_
    for <http{,s}://{,www.}code{-golf.io,.golf}>;

done_testing;
