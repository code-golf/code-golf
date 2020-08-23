use t;

constant $path = '/foo?bar=baz';

# TODO
# is $ua->get( $_ . PATH )->{url}, 'https://code.golf' . PATH, $_
#     for <http{,s}://{,www.}code{-golf.io,.golf}>;

done-testing;
