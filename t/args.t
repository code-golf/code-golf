use feature 'state';

use HTTP::Tiny;
use JSON::PP;
use Test2::V0;

local $/ = "\n\n";

for ( map [ split /\n/, $_, 2 ], <DATA> ) {
    my %run;    # Pick a hole that will definitely have unicode.
    @run{qw/Lang Code Hole/} = ( @$_, 'rock-paper-scissors-spock-lizard' );

    my $res = ( state $ua = HTTP::Tiny->new )->post(
        'https://code.golf/solution', { content => encode_json \%run } );

    die $res->{content} unless $res->{success};

    $res = decode_json $res->{content};

    is $res->{Out}, join( "\n", $res->{Argv}->@* ), $_->[0];
}

done_testing;

# haskell, lua, python, rust, and swift all have - as the first arg :-(
# https://rosettacode.org/wiki/Command-line_arguments

# Note that there are two ways to write an F# program - with and without an EntryPoint.
# Make sure they both work. The version without an EntryPoint is more useful when there
# are no arguments. Ideally it wouldn't be necessary to skip two arguments in that version,
# but it doesn't really matter since the other form should be used.
__DATA__
bash
for a; do echo $a; done

brainfuck
++++++++++>,[[.,]<.>,]

c
#include <stdio.h>
i; main(int n, char **a) { while(++i < n) puts(a[i]); }

c-sharp
class A {static void Main(string[] args){foreach(var a in args)System.Console.WriteLine(a);}}

f-sharp
[<EntryPoint>]
let main args =
    args |> Array.iter (printfn "%s")
    0

f-sharp
System.Environment.GetCommandLineArgs() |> Array.skip 2 |> Array.iter (printfn "%s")

fortran
character(10)::a
do i=1,iargc()
call getarg(i,a)
write(*,'(a)')a
enddo
end

go
package main
import "fmt"
import "os"
func main() { for _, a := range os.Args[1:] { fmt.Println(a) } }

haskell
import System.Environment;main=do x<-getArgs;mapM putStrLn$drop 1 x

j
echo>2}.ARGV

java
class A{public static void main(String[] args){for (var s:args){System.out.println(s);}}}

javascript
arguments.forEach(a => print(a))

julia
for a in ARGS; println(a); end

lisp
(dolist(x *args*)(format t "~A~&" x))

lua
for i=1, #arg do
    print(arg[i])
end

nim
import os
for a in commandLineParams(): echo a

perl
say for @ARGV

php
while($a = next($argv)) echo "$a\n"

powershell
$args

python
import sys
[print(a) for a in sys.argv[1:]]

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

swift
for a in CommandLine.arguments[1...] {
    print(a)
}
