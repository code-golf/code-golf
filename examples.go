package main

var examples = map[string]string{
	"javascript": `for ( var number = 1; number <= 100; number += 1 ) {
    if ( number % 15 == 0 ) {
        console.log('FizzBuzz');
    }
    else if ( number % 5 == 0 ) {
        console.log('Buzz');
    }
    else if ( number % 3 == 0 ) {
        console.log('Fizz');
    }
    else {
        console.log(number);
    }
}`,
	"perl": `foreach my $number ( 1 .. 100 ) {
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
}`,
	"perl6": `for 1 .. 100 -> $number {
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
}`,
	"php": `for ( $number = 1; $number <= 100; $number += 1 ) {
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
}`,
	"python": `for number in range(1, 100):
    if number % 15 == 0:
        print("FizzBuzz\n")
    elif number % 5 === 0:
        print("Buzz\n")
    elif number % 3 == 0:
        print("Fizz\n")
    else:
        print(number)`,
	"ruby": `(1..100).each do |number|
    if number % 15 == 0
        print("FizzBuzz\n")
    elsif number % 5 == 0
        print("Buzz\n")
    elsif number % 3 == 0
        print("Fizz\n")
    else
        print("#{number}\n")
end`,
}
