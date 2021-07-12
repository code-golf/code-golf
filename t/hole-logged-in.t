# This file contains tests for the hole route where being logged in is an important characteristic.
# Some of these tests have similar counterparts in hole-logged-out.t.

use hole;
use t;

plan 12;

sub createUserAndSession {
    my $dbh = dbh;
    my $userId = 1;
    createUser($dbh, $userId);
    my $session = createSession($dbh, $userId);
}

subtest 'Successful solutions are loaded from the database on reload.' => {
    plan 3;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: $raku57_55;
    $wd.isBytesAndChars: 57, 55;
    $wd.run;
    $wd.isPassing;
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isBytesAndChars: 57, 55, 'The byte count should be the same after reloading the page.';
}

subtest 'Untested solutions are loaded from localStorage on reload.' => {
    plan 3;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 3, 3, 'The byte count should be the same after reloading the page.';
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 0, 0, 'Untested solutions should not be preserved, after clearing localStorage.';
}

subtest 'Failing solutions are loaded from localStorage on reload.' => {
    plan 4;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.run;
    $wd.isFailing;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 3, 3, 'The byte count should be the same after reloading the page.';
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 0, 0, 'Failing solutions should not be preserved, after clearing localStorage.';
}

subtest 'The solution picker appears automatically, switching to bytes, and is independent of the scoring.' => {
    plan 8;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, '', "The solution picker shouldn't be visible, because a single solution optimizes both metrics.";
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.typeCode: BACKSPACE x 88 ~ $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should be active, because the submitted solution optimized that criteria.';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
}

subtest 'The solution picker appears automatically, switching to chars, and is independent of the scoring.' => {
    plan 8;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, '', "The solution picker shouldn't be visible, because a single solution optimizes both metrics.";
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, because the submitted solution optimized that criteria.';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
}

subtest 'The user can choose from bytes and chars solutions, independently of the scoring.' => {
    plan 13;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, because the submitted solution optimized that criteria.';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.isBytesAndChars: 210, 88;
    $wd.setSolution: 'bytes';
    is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should still be active, after reloading the page.';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.isBytesAndChars: 121, 121;
    $wd.setSolution: 'chars';
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, because the submitted solution optimized that criteria.';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.isBytesAndChars: 210, 88;
}

subtest 'The solution picker disappears automatically.' => {
    plan 8;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, because the submitted solution optimized that criteria.';
    $wd.typeCode: BACKSPACE x 88 ~ $python62_62;
    $wd.isBytesAndChars: 62, 62;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, '', "The solution picker shouldn't be visible, because a single solution once again optimizes both metrics.";
}

subtest 'Different bytes and chars solutions, and the active solution, are loaded from the database on reload.' => {
    plan 15;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    $wd.find('textarea').send-keys: BACKSPACE x 121;
    $wd.typeCode: $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, because the submitted solution optimized that criteria.';
    $wd.loadFizzBuzz;
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should still be active, after reloading the page.';
    $wd.isBytesAndChars: 210, 88;
    $wd.setSolution: 'bytes';
    is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should be active, after selecting it.';
    $wd.isBytesAndChars: 121, 121;
    $wd.loadFizzBuzz;
    is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should still be active, after reloading the page.';
    $wd.isBytesAndChars: 121, 121;
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should be active, by default, after clearing localStorage and reloading the page.';
    $wd.isBytesAndChars: 121, 121;
    $wd.setSolution: 'chars';
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, after selecting it.';
    $wd.isBytesAndChars: 210, 88;
}

subtest 'Different bytes and chars solutions, submitted while not logged in, can be resubmitted once logged in.' => {
    plan 14;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz; 
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    $wd.find('textarea').send-keys: BACKSPACE x 121;
    $wd.typeCode: $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    # Log in and reload the page.
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 210, 88;
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, because the submitted solution optimized that criteria.';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing;
    $wd.setSolution: 'bytes';
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.accept-alert;
    $wd.isBytesAndChars: 121, 121;
    # TODO: the solution picker should show bytes, but it actually disappears at this point, which is technically a bug.
    # is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should be active, after selecting it.';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing;
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    is $wd.getSolutionPickerState, 'bytes', 'The bytes solution should be active, by default, after clearing localStorage and reloading the page.';
    $wd.isBytesAndChars: 121, 121;
    $wd.setSolution: 'chars';
    is $wd.getSolutionPickerState, 'chars', 'The chars solution should be active, after selecting it.';
    $wd.isBytesAndChars: 210, 88;
}

subtest 'Users can choose not to restore autosaved solutions.' => {
    plan 6;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: $raku57_55;
    $wd.isBytesAndChars: 57, 55;
    $wd.run;
    $wd.isPassing;
    $wd.typeCode: BACKSPACE x 55 ~ 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.run;
    $wd.isFailing;
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.dismiss-alert;
    $wd.isBytesAndChars: 57, 55;
}

subtest 'Users can restore autosaved solutions.' => {
    plan 6;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: $raku57_55;
    $wd.isBytesAndChars: 57, 55;
    $wd.run;
    $wd.isPassing;
    $wd.typeCode: BACKSPACE x 55 ~ 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.run;
    $wd.isFailing;
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.accept-alert;
    $wd.isBytesAndChars: 3, 3;
}

subtest 'If the user improves their solution on another browser, they are not prompted to restore their old one.' => {
    plan 4;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    my $session = createUserAndSession;
    $wd.setSessionCookie: $session;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    # Improve the solution outside of this browser session.
    ok post-solution(:code($python62_62), :hole<fizz-buzz>, :lang<python>, :$session)<Pass>, 'Passes';
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 62, 62, 'The byte count should be lower, after reloading the page.';
}
