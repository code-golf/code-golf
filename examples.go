package main

var examples = map[string]string{
	"javascript": `for ( var number = 1; number <= 100; number += 1 ) {
    if ( number % 15 == 0 ) {
        print('FizzBuzz');
    }
    else if ( number % 5 == 0 ) {
        print('Buzz');
    }
    else if ( number % 3 == 0 ) {
        print('Fizz');
    }
    else {
        print(number);
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
    elseif ( $number % 5 == 0 ) {
        print("Buzz\n");
    }
    elseif ( $number % 3 == 0 ) {
        print("Fizz\n");
    }
    else {
        print("$number\n");
    }
}`,
	"python": `for number in range(1, 100):
    if number % 15 == 0:
        print("FizzBuzz")
    elif number % 5 == 0:
        print("Buzz")
    elif number % 3 == 0:
        print("Fizz")
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
    end
end`,
}
