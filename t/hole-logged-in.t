# This file contains tests for the hole route where being logged in is an important characteristic.
# Some of these tests have similar counterparts in hole-logged-out.t.

use hole;
use t;

plan 21;

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

subtest 'After manually reverting unsaved changes, the restore solution link is not shown.' => {
    plan 9;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
    $wd.typeCode: $raku57_55;
    $wd.isBytesAndChars: 57, 55, 'after typing code.';
    $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
    $wd.run;
    $wd.isPassing;
    $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
    $wd.typeCode: 'abc';
    $wd.isBytesAndChars: 60, 58, 'after modifying code.';
    $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
    $wd.typeCode: BACKSPACE x 3;
    $wd.isBytesAndChars: 57, 55, 'after manually reverting changes.';
    $wd.isRestoreSolutionLinkVisible: False, 'after manually reverting changes.';
}

subtest 'After submitting a shorter solution, the restore solution link is not shown.' => {
    plan 9;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
    $wd.typeCode: $raku59_57;
    $wd.isBytesAndChars: 59, 57, 'after typing code.';
    $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
    $wd.run;
    $wd.isPassing;
    $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
    $wd.typeCode: BACKSPACE x 57 ~ $raku57_55;
    $wd.isBytesAndChars: 57, 55, 'after modifying code.';
    $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
    $wd.run;
    $wd.isPassing;
    $wd.isRestoreSolutionLinkVisible: False, 'after submitting a shorter solution.';
}

for (False, True) -> $reloadFirst {
    my $context = $reloadFirst ?? ' and reloading the page' !! '';

    subtest "After submitting a longer solution$context, the restore solution link is shown and it works." => {
        plan 12;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        $wd.loadFizzBuzz;
        $wd.setSessionCookie: createUserAndSession;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Raku').click;
        $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
        $wd.typeCode: $raku57_55;
        $wd.isBytesAndChars: 57, 55, 'after typing code.';
        $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
        $wd.run;
        $wd.isPassing;
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
        $wd.typeCode: BACKSPACE x 55 ~ $raku59_57;
        $wd.isBytesAndChars: 59, 57, 'after modifying code.';
        $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
        $wd.run;
        $wd.isPassing;
        $wd.loadFizzBuzz if $reloadFirst;
        $wd.isRestoreSolutionLinkVisible: True, 'after submitting a longer solution.';
        $wd.restoreSolution;
        $wd.isBytesAndChars: 57, 55, 'after typing code.';
        $wd.isRestoreSolutionLinkVisible: False, 'after restoring solution';
        $wd.run;
        $wd.isPassing;
    }

    subtest "Successful solution can be restored after typing an untested solution$context." => {
        plan 17;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        $wd.loadFizzBuzz;
        $wd.setSessionCookie: createUserAndSession;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Raku').click;
        $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
        $wd.typeCode: $raku57_55;
        $wd.isBytesAndChars: 57, 55, 'after typing code.';
        $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
        $wd.isSolutionPickerState: '', 'before submitting code.';
        $wd.run;
        $wd.isPassing;
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
        $wd.isSolutionPickerState: '', 'after submitting code.';
        $wd.typeCode: 'abc';
        $wd.loadFizzBuzz if $reloadFirst;
        $wd.isBytesAndChars: 60, 58, 'after modifying code.';
        $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
        $wd.isSolutionPickerState: '', 'after modifying code.';
        $wd.restoreSolution;
        $wd.isBytesAndChars: 57, 55, 'after restoring the solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after restoring the solution.';
        $wd.isSolutionPickerState: '', 'after restoring the solution.';
        $wd.run;
        $wd.isPassing;
        # Reload the page to verify that the autosaved solution is gone.
        $wd.loadFizzBuzz;
        $wd.isBytesAndChars: 57, 55, 'after reloading the page.';
        $wd.isRestoreSolutionLinkVisible: False, 'after reloading the page.';
        $wd.isSolutionPickerState: '', 'after reloading the page.';
    }

    subtest "Successful solution can be restored after submitting a failed solution$context." => {
        plan 15;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        $wd.loadFizzBuzz;
        $wd.setSessionCookie: createUserAndSession;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Raku').click;
        $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
        $wd.typeCode: $raku57_55;
        $wd.isBytesAndChars: 57, 55, 'after typing code.';
        $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
        $wd.run;
        $wd.isPassing;
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
        $wd.typeCode: 'abc';
        $wd.isBytesAndChars: 60, 58, 'after modifying code.';
        $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
        $wd.run;
        $wd.isFailing;
        $wd.loadFizzBuzz if $reloadFirst;
        $wd.isRestoreSolutionLinkVisible: True, 'after submitting failing solution.';
        $wd.restoreSolution;
        $wd.isBytesAndChars: 57, 55, 'after restoring the solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after restoring the solution.';
        $wd.run;
        $wd.isPassing;
        # Reload the page to verify that the autosaved solution is gone.
        $wd.loadFizzBuzz;
        $wd.isBytesAndChars: 57, 55, 'after reloading the page.';
        $wd.isRestoreSolutionLinkVisible: False, 'after reloading the page.';
        $wd.isSolutionPickerState: '', 'after reloading the page.';
    }
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

for (False, True) -> $reloadFirst {
    my $context = $reloadFirst ?? ', and reloading the page' !! '';

    subtest "Successful solutions for both bytes and chars can be restored after typing untested solutions$context." => {
        plan 25;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        $wd.loadFizzBuzz;
        $wd.setSessionCookie: createUserAndSession;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Python').click;
        $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
        # Submit different solutions for bytes and chars.
        $wd.typeCode: $python121_121;
        $wd.isRestoreSolutionLinkVisible: False, 'after typing the bytes solution.';
        $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
        $wd.run;
        $wd.isPassing;
        $wd.isSolutionPickerState: '', 'after submitting the bytes solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting the bytes solution.';
        $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
        $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
        $wd.isRestoreSolutionLinkVisible: True, 'after typing the chars solution.';
        $wd.run;
        $wd.isPassing;
        $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting the chars solution.';
        # Modify both solutions.
        $wd.typeCode: 'A';
        $wd.isBytesAndChars: 211, 89, 'after modifying the chars solution.';
        $wd.isRestoreSolutionLinkVisible: True, 'after modifying the chars solution.';
        $wd.setSolution: 'bytes';
        $wd.isBytesAndChars: 121, 121, 'after switching to the bytes solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after after switching to the bytes solution.';
        $wd.typeCode: 'A';
        # Restore the solutions.
        $wd.loadFizzBuzz if $reloadFirst;
        $wd.isBytesAndChars: 122, 122, 'after modifying the bytes solution.';
        $wd.isRestoreSolutionLinkVisible: True, 'after modifying the bytes solution.';
        $wd.restoreSolution;
        $wd.isBytesAndChars: 121, 121, 'after restoring the bytes solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after restoring the bytes solution.';
        $wd.setSolution: 'chars';
        $wd.isBytesAndChars: 211, 89, 'after switching to the chars solution.';
        $wd.isRestoreSolutionLinkVisible: True, 'after switching to the chars solution.';
        $wd.restoreSolution;
        $wd.isBytesAndChars: 210, 88, 'after restoring the bytes solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after restoring the bytes solution.';
        $wd.setSolution: 'bytes';
        $wd.isBytesAndChars: 121, 121, 'after switching to the bytes solution again.';
        $wd.isRestoreSolutionLinkVisible: False, 'after after switching to the bytes solution again.';
    }
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
    plan 10;
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
    $wd.restoreSolution;
    $wd.isBytesAndChars: 210, 88, 'after discarding the bytes solution.';
    $wd.isSolutionPickerState: '', 'after discarding the bytes solution.';
}

subtest 'If the user improves their solution on another browser, the restore solution link is not shown.' => {
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
    # Improve the solution outside of this browser session.
    ok post-solution(:code($python62_62), :hole<fizz-buzz>, :lang<python>, :$session)<Pass>, 'Passes';
    $wd.loadFizzBuzz;
    $wd.isRestoreSolutionLinkVisible: False, 'after improving the solution outside of the browser session and reloading.';
    $wd.isBytesAndChars: 62, 62, 'The byte count should be lower, after reloading the page.';
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}
