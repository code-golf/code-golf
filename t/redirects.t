use HTTP::Tiny;
use Test2::V0;

use constant PATH => '/foo?bar=baz';

my $ua = HTTP::Tiny->new;

is $ua->get( $_ . PATH )->{url}, 'https://code-golf.io' . PATH, $_
    for <http{,s}://{,www.}code-golf.io>;

done_testing;
