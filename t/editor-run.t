use Test;
use WebDriver;

my $wd = WebDriver.new: :4444port, :host<firefox>,
    :capabilities(:alwaysMatch('moz:firefoxOptions' => :args('-headless',)));

$wd.get: 'https://code.golf/fizz-buzz';

$wd.find('Raku', :using(WebDriver::Selector::LinkText)).click;

$wd.find('textarea').send-keys:
    'say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1…100';

is $wd.find('#chars').text, '57 bytes, 55 chars';

$wd.find('Run', :using(WebDriver::Selector::LinkText)).click;

for ^5 {
    if my $text = $wd.find('h2').text {
        is $text, 'Pass 😀';
        last;
    }

    sleep 1;
}

done-testing;
