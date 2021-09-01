use t;

my $code = q:to/JS/;
for (let i = 1; i <= 100; i++)
    if (i % 15 == 0)
        print('FizzBuzz');
    else if (i % 3 == 0)
        print('Fizz');
    else if (i % 5 == 0)
        print('Buzz');
    else
        print(i);
JS

# -1 for the null termination.
$code ~= ' ' x 128 * 1024 - $code.chars - 1;

ok post-solution(:$code, :lang<javascript>)<Pass>, '128 KiB passes';

$code ~= ' ';

throws-like { post-solution :$code, :lang<javascript> }, Exception,
    :message(/"413 Request Entity Too Large"/), '128 KiB + 1 fails';

done-testing;
