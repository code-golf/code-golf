use feature 'state';

use HTTP::Tiny;
use JSON::PP;
use Test2::V0;

subtest $_->[0] => sub {
    my %run;    # Pick a hole that will definitely have unicode.
    @run{qw/Lang Code Hole/} = ( @$_, 'rock-paper-scissors-spock-lizard' );

    my $res = ( state $ua = HTTP::Tiny->new )->post(
        'https://code-golf.io/solution', { content => encode_json \%run } );

    die $res->{content} unless $res->{success};

    $res = decode_json $res->{content};

    is $res->{Err}, '', 'Err';
    is $res->{Out}, join( "\n", $res->{Argv}->@* ), 'Out';
} for do { local $/ = "\n\n"; map [ split /\n/, $_, 2 ], <DATA> };

done_testing;

# lua, python, rust all have - as the first arg :-(
# https://rosettacode.org/wiki/Command-line_arguments
__DATA__
bash
for a; do echo $a; done

c
#include <stdio.h>
i; main(int n, char **a) { while(++i < n) puts(a[i]); }

javascript
arguments.forEach(a => print(a))

julia
for a in ARGS; println(a); end

lua
for i=1, #arg do
    print(arg[i])
end

nim
import os
for i in 1 .. paramCount(): echo paramStr(i)

nim
import os
for a in commandLineParams(): echo a

perl
say for @ARGV

php
while($a = next($argv)) echo "$a\n"

python
import sys
[print(a)for a in sys.argv[1:]]

raku
@*ARGS».say

ruby
puts ARGV

rust
use std::env;
fn main() {
    for arg in env::args() {
        if arg != "-" {
            println!("{}", arg);
        }
    }
}
