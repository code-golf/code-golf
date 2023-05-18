use Test;
use WebDriver;

# https://www.w3.org/TR/webdriver/#keyboard-actions
constant WD-BACKSPACE = "\xe003";
constant WD-CONTROL   = "\xe009";
constant WD-END       = "\xe010";
constant WD-DELETE    = "\xe017";

our $raku57_55 is export = 'say "Fizz" x $_ %% 3 ~ "Buzz" x $_ %% 5 || $_ for 1…100';

our $raku59_57 is export = 'say ("Fizz" x $_ %% 3) ~ "Buzz" x $_ %% 5 || $_ for 1…100';

our $python121_121 is export = q:to/CODE/.trim;
    for x in range(1,101):
     if x%15==0:print('FizzBuzz')
     elif x%3==0:print('Fizz')
     elif x%5==0:print('Buzz')
     else:print(x)
    CODE

our $python210_88 is export = "exec('景爠砠楮⁲慮来⠱ⰱ〱⤺ਠ楦⁸┱㔽㴰㩰物湴⠧䙩空䉵空✩ਠ敬楦⁸┳㴽〺灲楮琨❆楺稧⤊⁥汩映砥㔽㴰㩰物湴⠧䉵空✩ਠ敬獥㩰物湴⡸⤠'.encode('utf-16be'))";

our $python62_62 is export = 'for i in range(1,101):print("Fizz"*(i%3<1)+"Buzz"*(i%5<1)or i)';

our $python_fibonacci_66_66 is export = q:to/CODE/.trim;
    x = 0
    y = 1
    for i in range(31):
     print(x)
     z = x + y
     x = y
     y = z
    CODE

our $python_fibonacci_126_60 is export = "exec('砽《礽ㄊ景爠椠楮⁲慮来⠳ㄩ㨊††灲楮琨砩ਠ†⁺㵸⭹ਠ†⁸㵹ਠ†⁹㵺'.encode('utf-16be'))";

class HoleWebDriver is WebDriver is export {
    method create(::?CLASS:U $wd:) {
        $wd.new: :4444port, :host<firefox>, :capabilities(:alwaysMatch(
            {:acceptInsecureCerts, 'moz:firefoxOptions' => {:args('-headless',), :prefs('devtools.console.stdout.content' => True)}}));
    }

    method clearLocalStorage {
        $.js: 'localStorage.clear();';
    }

    method findAndWait(Str:D $text, WebDriver::Selector:D :$using = CSS) {
        for ^5 {
            try {
                return $.find($text, :$using);
            }
            sleep 1;
        }

        $.find($text, :$using);
    }

    method getLangLink(Str:D $lang) {
        $.findAndWait: $lang, :using(PartialLinkText);
    }

    method getLanguageActive(Str:D $lang) {
        $.getLangLink($lang).prop('href') eq '';
    }

    method getSolutionLink(Str:D $solution) {
        # Don't try again, if not found. Sometimes these links aren't present.
        $.find("Fewest $solution.tc()", :using(PartialLinkText)).first(*.visible);
    }

    method getScoringLink(Str:D $scoring) {
        $.findAndWait($scoring.tc, :using(LinkText));
    }

    method getSolutionPickerState of Str {
        try {
            return "bytes" if $.getSolutionLink('bytes').prop('href') eq '';
            return "chars" if $.getSolutionLink('chars').prop('href') eq '';
        }
        '';
    }

    method getScoringPickerState of Str {
        try {
            return "bytes" if $.getScoringLink('bytes').prop('href') eq '';
            return "chars" if $.getScoringLink('chars').prop('href') eq '';
        }
        '';
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isBytesAndChars(Int:D $bytes, Int:D $chars, Str:D $context) {
        is $.find("#strokes").text, "$bytes bytes, $chars chars",
            "Confirm byte and char counts, $context";
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isFailing(Str:D $context) {
        $.isResult: 'Fail ☹️', $context;
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isPassing(Str:D $context) {
        $.isResult: 'Pass 😀', $context;
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isResult(Str:D $expectedText, Str:D $context) {
        for ^5 {
            if (my $text = $.find('h2').text) && $text ne '…' {
                is $text, $expectedText, "Confirm the result of running the program, $context";
                return;
            }

            sleep 1;
        }

        flunk "Failed to find run results, $context";
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isScoringPickerState(Str:D $expectedState, Str:D $context) {
        is $.getScoringPickerState, $expectedState, "Confirm the scoring picker state, $context";
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isSolutionPickerState(Str:D $expectedState, Str:D $context) {
        my $desc = $expectedState ??
            "The $expectedState solution should be active" !!
            "The solution picker shouldn't be visible, because a single solution optimizes both metrics";
        is $.getSolutionPickerState, $expectedState, "$desc, $context";
    }

    # Methods whose names begin with "is" do exactly one assertion.
    method isRestoreSolutionLinkVisible(Bool:D $expectedState, Str:D $context) {
        my $callable = { $.find('Restore solution', :using(WebDriver::Selector::LinkText)) };
        if $expectedState {
            lives-ok $callable, "Restore solution link should be visible, $context";
        }
        else {
            throws-like $callable, Exception, :message(/'Unable to locate element'/),
                "Restore solution link should be hidden, $context";
        }
    }

    method loadFibonacci { $.get: 'https://app/fibonacci' }
    method loadFizzBuzz  { $.get: 'https://app/fizz-buzz' }

    method restoreSolution {
        $.find('Restore solution', :using(LinkText)).click;
    }

    method run {
        $.find("#runBtn").click;
    }

    method setScoring(Str:D $scoring) {
        $.getScoringLink($scoring).click;
    }

    method setSolution(Str:D $solution) {
        $.getSolutionLink($solution).click;
    }

    method setSessionCookie(Str:D $session) {
        $.cookie: '__Host-session', $session, :httpOnly, :sameSite<Lax>, :secure;
    }

    method clearCode {
        $.find('.cm-content').click.send-keys: WD-CONTROL ~ 'a' ~ WD-DELETE;
    }

    method typeCode(Str:D $code) {
        $.find('.cm-content').click.send-keys(WD-CONTROL ~ WD-END).send-keys($code);
    }
}
