use HTTP::Tiny;
use Test2::V0;

use constant PATH => '/foo?bar=baz';

my $ua = HTTP::Tiny->new( max_redirect => 0 );

for (qw(
          http://86.8.162.46
         http://code-golf.io
     http://www.code-golf.io
    https://www.code-golf.io
)) {
    like $ua->get( $_ . PATH ) => {
        headers => { location => 'https://code-golf.io' . PATH },
        status  => 308,
    } => $_;
}

done_testing;
