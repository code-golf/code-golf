package main

const fizzBuzzExample = `foreach my $number ( 1 .. 100 ) {
    if ( $number % 15 == 0 ) {
        print("FizzBuzz\n");
    }
    elsif ( $number % 5 == 0 ) {
        print("Buzz\n");
    }
    elsif ( $number % 3 == 0 ) {
        print("Fizz\n");
    }
    else {
        print("$number\n");
    }
}`
