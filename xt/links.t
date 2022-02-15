use HTTP::Tiny;
use Test;
use TOML::Thumb;

my $ua = HTTP::Tiny.new :!max-redirect;

# Ensure all the links in the holes resolve with no redirects.
for 'config/holes.toml'.IO.&from-toml.values {
    next unless .<links>;

    for .<links>».<url> {
        s|^ '//' |https://|;
        s|  "'"  |%27|;

        is $ua.head($_)<status>, 200, $_;
    };
}

done-testing;
