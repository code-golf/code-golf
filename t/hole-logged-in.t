# This file contains tests for the hole route where being logged in is an important characteristic.
# Some of these tests have similar counterparts in hole-logged-out.t.

use hole;
use t;

plan 17;

sub createUserAndSession {
    my $dbh = dbh;
    my $userId = 1;
    createUser($dbh, $userId);
    my $session = createSession($dbh, $userId);
}

subtest 'Successful solutions are loaded from the database on reload.' => {
    plan 5;
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
    $wd.isSolutionPickerState: '';
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isBytesAndChars: 57, 55, 'after clearing localStorage and reloading the page.';
    $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
}

subtest 'Untested solutions are loaded from localStorage on reload.' => {
    plan 5;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 3, 3, 'after reloading the page.';
    $wd.isSolutionPickerState: '';
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isBytesAndChars: 0, 0, 'after clearing localStorage and reloading the page.';
    $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
}

subtest 'Failing solutions are loaded from localStorage on reload.' => {
    plan 5;
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
    $wd.isBytesAndChars: 3, 3, 'after reloading the page.';
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isBytesAndChars: 0, 0, 'after clearing localStorage and reloading the page.';
    $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
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
    $wd.isSolutionPickerState: '';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.typeCode: BACKSPACE x 88 ~ $python121_121;
    $wd.isBytesAndChars: 121, 121;
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: 'bytes';
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
    $wd.isSolutionPickerState: '';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88;
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: 'chars';
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
    $wd.isSolutionPickerState: 'chars';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.isBytesAndChars: 210, 88;
    $wd.setSolution: 'bytes';
    $wd.isSolutionPickerState: 'bytes', 'after reloading the page.';
    is $wd.getScoringPickerState, 'bytes', "The scoring should be the default, bytes.";
    $wd.isBytesAndChars: 121, 121;
    $wd.setSolution: 'chars';
    $wd.isSolutionPickerState: 'chars', 'after switching solutions.';
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
    $wd.isSolutionPickerState: 'chars';
    $wd.typeCode: BACKSPACE x 88 ~ $python62_62;
    $wd.isBytesAndChars: 62, 62;
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: '';
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
    $wd.isSolutionPickerState: 'chars';
    $wd.loadFizzBuzz;
    $wd.isSolutionPickerState: 'chars', 'after reloading the page.';
    $wd.isBytesAndChars: 210, 88;
    $wd.setSolution: 'bytes';
    $wd.isSolutionPickerState: 'bytes', 'after switching solutions.';
    $wd.isBytesAndChars: 121, 121;
    $wd.loadFizzBuzz;
    $wd.isSolutionPickerState: 'bytes', 'after reloading the page.';
    $wd.isBytesAndChars: 121, 121;
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.isSolutionPickerState: 'bytes', 'after clearing localStorage and reloading the page.';
    $wd.isBytesAndChars: 121, 121;
    $wd.setSolution: 'chars';
    $wd.isSolutionPickerState: 'chars', 'after switching solutions.';
    $wd.isBytesAndChars: 210, 88;
}

# TODO: Add a variant that submits the two solutions in the opposite order.
subtest 'After submitting different bytes and chars solutions while not logged in, users can submit both once logged in.' => {
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
    $wd.isSolutionPickerState: 'chars';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing;
    $wd.setSolution: 'bytes';
    $wd.isBytesAndChars: 121, 121;
    is $wd.getSolutionPickerState, '', 'after submitting the first solution after logging in.';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing;
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.isSolutionPickerState: 'bytes', 'after clearing localStorage and reloading the page.';
    $wd.isBytesAndChars: 121, 121;
    $wd.setSolution: 'chars';
    $wd.isSolutionPickerState: 'chars', 'after switching solutions.';
    $wd.isBytesAndChars: 210, 88;
}

subtest 'After submitting different bytes and chars solutions while not logged in, users can submit one and discard the other once logged in.' => {
    plan 13;
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
    $wd.isSolutionPickerState: 'chars';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing;
    $wd.setSolution: 'bytes';
    $wd.isBytesAndChars: 121, 121;
    # Reload the page, without submitting the solution.
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.dismiss-alert;
    $wd.isBytesAndChars: 210, 88, 'after reloading the page.';
    is $wd.getSolutionPickerState, '', 'after reloading the page.';
    # Reload the page to verify that users aren't prompted to restore the solution again.
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 210, 88, 'after reloading the page again.';
    is $wd.getSolutionPickerState, '', 'after reloading the page again.';
}

subtest 'Users can choose not to restore autosaved solutions.' => {
    plan 11;
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
    $wd.isSolutionPickerState: '';
    $wd.typeCode: BACKSPACE x 55 ~ 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.run;
    $wd.isFailing;
    $wd.isSolutionPickerState: '', 'after submitting failing solution.';
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.dismiss-alert;
    $wd.isBytesAndChars: 57, 55;
    $wd.isSolutionPickerState: '', 'after reloading the page.';
    # Reload the page to verify that users aren't prompted to restore the solution again.
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 57, 55;
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}

subtest 'Users can restore autosaved solutions.' => {
    plan 9;
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
    $wd.isSolutionPickerState: '';
    $wd.typeCode: BACKSPACE x 55 ~ 'abc';
    $wd.isBytesAndChars: 3, 3;
    $wd.run;
    $wd.isFailing;
    $wd.isSolutionPickerState: '', 'after submitting failing solution.';
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.accept-alert;
    $wd.isBytesAndChars: 3, 3;
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}

subtest 'Discarding autosaved solutions applies to both bytes and chars solutions.' => {
    plan 15;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    # Submit different solutions for bytes and chars.
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: '';
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: 'chars';
    # Modify both solutions.
    $wd.typeCode: 'A';
    $wd.isBytesAndChars: 211, 89, 'after modifying the chars solution.';
    $wd.setSolution: 'bytes';
    $wd.typeCode: 'A';
    $wd.isBytesAndChars: 122, 122, 'after modifying the bytes solution.';
    # Reload the page and restore the local solutions.
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.dismiss-alert;
    $wd.isBytesAndChars: 121, 121, 'after reloading the page.';
    $wd.isSolutionPickerState: 'bytes', 'after reloading the page.';
    $wd.setSolution: 'chars';
    $wd.isBytesAndChars: 210, 88, 'after switching to the chars solution.';
    # Reload the page to verify that users aren't prompted to restore the solution again.
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 210, 88, 'after reloading the page again.';
    $wd.isSolutionPickerState: 'chars', 'after reloading the page again.';
    $wd.setSolution: 'bytes';
    $wd.isBytesAndChars: 121, 121, 'after switching to the bytes solution';
}

subtest 'Restoring autosaved solutions applies to both bytes and chars solutions.' => {
    plan 12;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    # Submit different solutions for bytes and chars.
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: '';
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: 'chars';
    # Modify both solutions.
    $wd.typeCode: 'A';
    $wd.isBytesAndChars: 211, 89, 'after modifying the chars solution.';
    $wd.setSolution: 'bytes';
    $wd.typeCode: 'A';
    $wd.isBytesAndChars: 122, 122, 'after modifying the bytes solution.';
    # Reload the page and restore the local solutions.
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.accept-alert;
    $wd.isBytesAndChars: 122, 122, 'after reloading the page.';
    $wd.isSolutionPickerState: 'bytes', 'after reloading the page.';
    $wd.setSolution: 'chars';
    $wd.isBytesAndChars: 211, 89, 'after switching to the chars solution.';
}

subtest 'If the user improves their solution on another browser, they are not prompted to restore their old one.' => {
    plan 6;
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
    $wd.isSolutionPickerState: '';
    # Improve the solution outside of this browser session.
    ok post-solution(:code($python62_62), :hole<fizz-buzz>, :lang<python>, :$session)<Pass>, 'Passes';
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 62, 62, 'The byte count should be lower, after reloading the page.';
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}

subtest 'Unsaved changes are auto-saved.' => {
    plan 7;
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
    $wd.isSolutionPickerState: '';
    # Type some code, but don't run it.
    $wd.typeCode: "a";
    $wd.isBytesAndChars: 122, 122;
    $wd.loadFizzBuzz;
    is $wd.alert-text, 'Your local copy of the code is different than the remote one. Do you want to restore the local version?', 'Confirm alert text';
    $wd.accept-alert;
    $wd.isBytesAndChars: 122, 122, 'The byte count should be the same, after reloading the page.';
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}

subtest 'If unsaved changes are manually reverted, the user is not prompted to restore them.' => {
    plan 7;
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
    $wd.isSolutionPickerState: '';
    # Type some code, but don't run it.
    $wd.typeCode: "a";
    $wd.isBytesAndChars: 122, 122;
    # Manually undo the code change.
    $wd.typeCode: BACKSPACE;
    $wd.isBytesAndChars: 121, 121;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 121, 121, 'The byte count should be the same, after reloading the page.';
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}
