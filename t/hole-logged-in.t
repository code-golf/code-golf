# This file contains tests for the hole route where the user's logged in state is relevant.

use hole;
use t;

plan 37;

sub createUserAndSession {
    my $dbh = dbh;
    my $userId = 1;
    createUser($dbh, $userId);
    my $session = createSession($dbh, $userId);
}

for (False, True) -> $loggedIn {
    my $loggedInContext = $loggedIn ?? 'logged in' !! 'not logged in';

    # When the user is logged in, some tests expect to be able to reload the page and observe the same behavior.
    # When the user is logged  out, they are not expected to work.
    my @reloadFirstValues = False;
    @reloadFirstValues.push: True if $loggedIn;

    sub setup(HoleWebDriver:D $wd) {
        $wd.loadFizzBuzz;
        if $loggedIn {
            $wd.setSessionCookie: createUserAndSession;
            $wd.loadFizzBuzz;
        }
    }

    subtest "When $loggedInContext, successful solutions persist on reload." => {
        plan 5;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Raku').click;
        $wd.typeCode: $raku57_55;
        $wd.isBytesAndChars: 57, 55, 'after typing code.';
        $wd.run;
        $wd.isPassing: 'after running code.';
        $wd.isSolutionPickerState: '', 'after running code.';
        $wd.clearLocalStorage;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Raku').click;
        $wd.isBytesAndChars: 57 * $loggedIn, 55 * $loggedIn, 'after clearing localStorage and reloading the page.';
        $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
    }

    subtest "When $loggedInContext, untested solutions are loaded from localStorage on reload." => {
        plan 5;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Raku').click;
        $wd.typeCode: 'abc';
        $wd.isBytesAndChars: 3, 3, 'after typing code.';
        $wd.loadFizzBuzz;
        $wd.isBytesAndChars: 3, 3, 'after reloading the page.';
        $wd.isSolutionPickerState: '', 'after reloading the page.';
        $wd.clearLocalStorage;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Raku').click;
        $wd.isBytesAndChars: 0, 0, 'after clearing localStorage and reloading the page.';
        $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
    }

    subtest "When $loggedInContext, failing solutions are loaded from localStorage on reload." => {
        plan 5;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Raku').click;
        $wd.typeCode: 'abc';
        $wd.isBytesAndChars: 3, 3, 'after typing a failing solution.';
        $wd.run;
        $wd.isFailing: 'after submitting a failing solution';
        $wd.loadFizzBuzz;
        $wd.isBytesAndChars: 3, 3, 'after reloading the page.';
        $wd.clearLocalStorage;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Raku').click;
        $wd.isBytesAndChars: 0, 0, 'after clearing localStorage and reloading the page.';
        $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
    }

    subtest "When $loggedInContext, after manually reverting unsaved changes, the restore solution link is not shown." => {
        plan 9;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Raku').click;
        $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
        $wd.typeCode: $raku57_55;
        $wd.isBytesAndChars: 57, 55, 'after typing code.';
        $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
        $wd.run;
        $wd.isPassing: 'after submitting code.';
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
        $wd.typeCode: 'abc';
        $wd.isBytesAndChars: 60, 58, 'after modifying code.';
        $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
        $wd.typeCode: BACKSPACE x 3;
        $wd.isBytesAndChars: 57, 55, 'after manually reverting changes.';
        $wd.isRestoreSolutionLinkVisible: False, 'after manually reverting changes.';
    }

    subtest "When $loggedInContext, after submitting a shorter solution, the restore solution link is not shown." => {
        plan 9;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Raku').click;
        $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
        $wd.typeCode: $raku59_57;
        $wd.isBytesAndChars: 59, 57, 'after typing code.';
        $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
        $wd.run;
        $wd.isPassing: 'after submitting code.';
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
        $wd.typeCode: BACKSPACE x 57 ~ $raku57_55;
        $wd.isBytesAndChars: 57, 55, 'after typing a shorter solution.';
        $wd.isRestoreSolutionLinkVisible: True, 'after typing a shorter solution.';
        $wd.run;
        $wd.isPassing: 'after submitting a shorter solution.';
        $wd.isRestoreSolutionLinkVisible: False, 'after submitting a shorter solution.';
    }

    for @reloadFirstValues -> $reloadFirst {
        my $context = $reloadFirst ?? ' and reloading the page' !! '';

        subtest "When $loggedInContext, after submitting a longer solution$context, the restore solution link is shown and it works." => {
            plan 12;
            my $wd = HoleWebDriver.create;
            LEAVE $wd.delete-session;
            setup $wd;
            $wd.getLangLink('Raku').click;
            $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
            $wd.typeCode: $raku57_55;
            $wd.isBytesAndChars: 57, 55, 'after typing code.';
            $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
            $wd.run;
            $wd.isPassing: 'after submitting code.';
            $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
            $wd.typeCode: BACKSPACE x 55 ~ $raku59_57;
            $wd.isBytesAndChars: 59, 57, 'after typing a longer solution.';
            $wd.isRestoreSolutionLinkVisible: True, 'after typing a longer solution.';
            $wd.run;
            $wd.isPassing: 'after submitting a longer solution.';
            $wd.loadFizzBuzz if $reloadFirst;
            $wd.isRestoreSolutionLinkVisible: True, 'after submitting a longer solution.';
            $wd.restoreSolution;
            $wd.isBytesAndChars: 57, 55, 'after typing code.';
            $wd.isRestoreSolutionLinkVisible: False, 'after restoring solution';
            $wd.run;
            $wd.isPassing: 'after restoring solution and running it.';
        }

        subtest "When $loggedInContext, a successful solution can be restored after typing an untested solution$context." => {
            plan 17;
            my $wd = HoleWebDriver.create;
            LEAVE $wd.delete-session;
            setup $wd;
            $wd.getLangLink('Raku').click;
            $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
            $wd.typeCode: $raku57_55;
            $wd.isBytesAndChars: 57, 55, 'after typing code.';
            $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
            $wd.isSolutionPickerState: '', 'before submitting code.';
            $wd.run;
            $wd.isPassing: 'after submitting code.';
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
            $wd.isPassing: 'after restoring solution and running it.';
            # Reload the page to verify that the autosaved solution is gone.
            $wd.loadFizzBuzz;
            $wd.isBytesAndChars: 57, 55, 'after reloading the page.';
            $wd.isRestoreSolutionLinkVisible: False, 'after reloading the page.';
            $wd.isSolutionPickerState: '', 'after reloading the page.';
        }

        subtest "When $loggedInContext, a successful solution can be restored after submitting a failed solution$context." => {
            plan 21;
            my $wd = HoleWebDriver.create;
            LEAVE $wd.delete-session;
            setup $wd;
            $wd.getLangLink('Raku').click;
            $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
            $wd.typeCode: $raku57_55;
            $wd.isBytesAndChars: 57, 55, 'after typing code.';
            $wd.isRestoreSolutionLinkVisible: False, 'before submitting code.';
            $wd.isSolutionPickerState: '', 'before submitting code.';
            $wd.run;
            $wd.isPassing: 'after submitting code.';
            $wd.isRestoreSolutionLinkVisible: False, 'after submitting code.';
            $wd.isSolutionPickerState: '', 'after submitting code.';
            $wd.typeCode: 'abc';
            $wd.isBytesAndChars: 60, 58, 'after modifying code.';
            $wd.isRestoreSolutionLinkVisible: True, 'after modifying code.';
            $wd.isSolutionPickerState: '', 'after modifying code.';
            $wd.run;
            $wd.isFailing: 'after submitting a failing solution';
            $wd.loadFizzBuzz if $reloadFirst;
            $wd.isRestoreSolutionLinkVisible: True, 'after submitting failing solution.';
            $wd.isSolutionPickerState: '', 'after submitting failing solution.';
            $wd.restoreSolution;
            $wd.isBytesAndChars: 57, 55, 'after restoring the solution.';
            $wd.isRestoreSolutionLinkVisible: False, 'after restoring the solution.';
            $wd.isSolutionPickerState: '', 'after restoring the solution.';
            $wd.run;
            $wd.isPassing: 'after restoring solution and running it.';
            $wd.isSolutionPickerState: '', 'after running the restored solution.';
            # Reload the page to verify that the autosaved solution is gone.
            $wd.loadFizzBuzz;
            $wd.isBytesAndChars: 57, 55, 'after reloading the page.';
            $wd.isRestoreSolutionLinkVisible: False, 'after reloading the page.';
            $wd.isSolutionPickerState: '', 'after reloading the page.';
        }
    }

    subtest "When $loggedInContext, the solution picker appears automatically, switching to bytes, and is independent of the scoring." => {
        plan 8;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Python').click;
        $wd.typeCode: $python210_88;
        $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the chars solution';
        $wd.isSolutionPickerState: '', 'after submitting the chars solution';
        $wd.isScoringPickerState: 'bytes', 'after submitting the chars solution. The scoring should be the default, bytes.';
        $wd.typeCode: BACKSPACE x 88 ~ $python121_121;
        $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the bytes solution.';
        $wd.isSolutionPickerState: 'bytes', 'after submitting the bytes solution.';
        $wd.isScoringPickerState: 'bytes', 'after submitting the bytes solution. The scoring should be the default, bytes.';
    }

    subtest "When $loggedInContext, the solution picker appears automatically, switching to chars, and is independent of the scoring." => {
        plan 8;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Python').click;
        $wd.typeCode: $python121_121;
        $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the bytes solution';
        $wd.isSolutionPickerState: '', 'after submitting the bytes solution';
        $wd.isScoringPickerState: 'bytes', 'after submitting the bytes solution. The scoring should be the default, bytes.';
        $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
        $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the chars solution';
        $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution';
        $wd.isScoringPickerState: 'bytes', 'after submitting the chars solution. The scoring should be the default, bytes.';
    }

    subtest "When $loggedInContext, the user can choose from bytes and chars solutions, independently of the scoring." => {
        plan 13;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Python').click;
        $wd.typeCode: $python121_121;
        $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the bytes solution.';
        $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
        $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the chars solution';
        $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution';
        $wd.isScoringPickerState: 'bytes', 'after submitting the chars solution. The scoring should be the default, bytes.';
        $wd.isBytesAndChars: 210, 88, 'after running code';
        $wd.setSolution: 'bytes';
        $wd.isBytesAndChars: 121, 121, 'after switching to the bytes solution.';
        $wd.isSolutionPickerState: 'bytes', 'after switching to the bytes solution.';
        $wd.isScoringPickerState: 'bytes', 'after switching to the bytes solution. The scoring should be the default, bytes.';
        $wd.setSolution: 'chars';
        $wd.isBytesAndChars: 210, 88, 'after switching to the chars solutions.';
        $wd.isSolutionPickerState: 'chars', 'after switching to the chars solutions.';
        $wd.isScoringPickerState: 'bytes', 'after switching to the chars solution. The scoring should be the default, bytes.';
    }

    subtest "When $loggedInContext, the solution picker disappears automatically." => {
        plan 8;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Python').click;
        $wd.typeCode: $python121_121;
        $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the bytes solution.';
        $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
        $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the chars solution';
        $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution';
        $wd.typeCode: BACKSPACE x 88 ~ $python62_62;
        $wd.isBytesAndChars: 62, 62, 'after typing a solution that improves both metrics.';
        $wd.run;
        $wd.isPassing: 'after running the solution that improves both metrics.';
        $wd.isSolutionPickerState: '', 'after running a solution that improves both metrics.';
    }

    subtest "When $loggedInContext, different bytes and chars solutions, and the active solution, persist on reload." => {
        plan 13 + 2 * $loggedIn;
        my $wd = HoleWebDriver.create;
        LEAVE $wd.delete-session;
        setup $wd;
        $wd.getLangLink('Python').click;
        $wd.typeCode: $python121_121;
        $wd.isBytesAndChars: 121, 121, 'after submitting the bytes solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the bytes solution.';
        $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
        $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
        $wd.run;
        $wd.isPassing: 'after submitting the chars solution';
        $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution';
        $wd.loadFizzBuzz;
        $wd.isSolutionPickerState: 'chars', 'after reloading the page.';
        $wd.isBytesAndChars: 210, 88, 'after reloading the page.';
        $wd.setSolution: 'bytes';
        $wd.isSolutionPickerState: 'bytes', 'after switching solutions.';
        $wd.isBytesAndChars: 121, 121, 'after switching solutions.';
        $wd.loadFizzBuzz;
        $wd.isSolutionPickerState: 'bytes', 'after reloading the page again.';
        $wd.isBytesAndChars: 121, 121, 'after reloading the page again.';
        $wd.clearLocalStorage;
        $wd.loadFizzBuzz;
        $wd.getLangLink('Python').click;
        if $loggedIn {
            $wd.isSolutionPickerState: 'bytes', 'after clearing localStorage and reloading the page.';
            $wd.isBytesAndChars: 121 * $loggedIn, 121 * $loggedIn, 'after clearing localStorage and reloading the page.';
            $wd.setSolution: 'chars';
            $wd.isSolutionPickerState: 'chars', 'after switching to the chars solution.';
            $wd.isBytesAndChars: 210 * $loggedIn, 88 * $loggedIn, 'after switching to the chars solution.';
        }
        else {
            $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
            $wd.isBytesAndChars: 0, 0, 'after clearing localStorage and reloading the page.';
        }
    }

    for @reloadFirstValues -> $reloadFirst {
        my $context = $reloadFirst ?? ', and reloading the page' !! '';

        subtest "When $loggedInContext, successful solutions for both bytes and chars can be restored after typing untested solutions$context." => {
            plan 25;
            my $wd = HoleWebDriver.create;
            LEAVE $wd.delete-session;
            setup $wd;
            $wd.getLangLink('Python').click;
            $wd.isRestoreSolutionLinkVisible: False, 'before typing code.';
            # Submit different solutions for bytes and chars.
            $wd.typeCode: $python121_121;
            $wd.isRestoreSolutionLinkVisible: False, 'after typing the bytes solution.';
            $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
            $wd.run;
            $wd.isPassing: 'after submitting the bytes solution.';
            $wd.isSolutionPickerState: '', 'after submitting the bytes solution.';
            $wd.isRestoreSolutionLinkVisible: False, 'after submitting the bytes solution.';
            $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
            $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
            $wd.isRestoreSolutionLinkVisible: True, 'after typing the chars solution.';
            $wd.run;
            $wd.isPassing: 'after submitting the chars solution.';
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
}

subtest 'When not logged in, after a user submits different bytes and chars solutions and then logs in, they can submit the chars solution and then the bytes solution.' => {
    plan 15;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged out.';
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged out.';
    $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution, while logged out.';
    # Log in and reload the page.
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 210, 88, 'after reloading the page';
    $wd.isSolutionPickerState: 'chars', 'after reloading the page';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged in.';
    $wd.setSolution: 'bytes';
    $wd.isBytesAndChars: 121, 121, 'after switching to the bytes solution.';
    $wd.isSolutionPickerState: '', 'after switching to the bytes solution. There is only one submitted solution at this point.';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged in.';
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.isSolutionPickerState: 'bytes', 'after clearing localStorage and reloading the page.';
    $wd.isBytesAndChars: 121, 121, 'after clearing localStorage and reloading the page.';
    $wd.setSolution: 'chars';
    $wd.isSolutionPickerState: 'chars', 'after switching to the chars solution.';
    $wd.isBytesAndChars: 210, 88, 'after switching to the chars solution.';
}

subtest 'When not logged in, after a user submits different bytes and chars solutions and then logs in, they can submit the bytes solution and then the chars solution.' => {
    plan 15;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python210_88;
    $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged out.';
    $wd.typeCode: BACKSPACE x 88 ~ $python121_121;
    $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged out.';
    $wd.isSolutionPickerState: 'bytes', 'after submitting the bytes solution, while logged out.';
    # Log in and reload the page.
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 121, 121, 'after reloading the page';
    $wd.isSolutionPickerState: 'bytes', 'after reloading the page';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged in.';
    $wd.setSolution: 'chars';
    $wd.isBytesAndChars: 210, 88, 'after switching to the chars solution.';
    $wd.isSolutionPickerState: '', 'after switching to the chars solution. There is only one submitted solution at this point.';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged in.';
    # Clear local storage just to prove that it's not required.
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.isSolutionPickerState: 'bytes', 'after clearing localStorage and reloading the page.';
    $wd.isBytesAndChars: 121, 121, 'after clearing localStorage and reloading the page.';
    $wd.setSolution: 'chars';
    $wd.isSolutionPickerState: 'chars', 'after switching to the chars solution.';
    $wd.isBytesAndChars: 210, 88, 'after switching to the chars solution.';
}

subtest 'When not logged in, after a user submits different bytes and chars solutions and then logs in, they can submit the chars solution and discard the bytes solution.' => {
    plan 12;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python121_121;
    $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged out.';
    $wd.typeCode: BACKSPACE x 121 ~ $python210_88;
    $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged out.';
    $wd.isSolutionPickerState: 'chars', 'after submitting the chars solution, while logged out.';
    # Log in and reload the page.
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 210, 88, 'after logging in and reloading the page.';
    $wd.isSolutionPickerState: 'chars', 'after logging in and reloading the page.';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged in.';
    $wd.setSolution: 'bytes';
    $wd.isBytesAndChars: 121, 121, 'after switching to the bytes solution.';
    $wd.isSolutionPickerState: '', 'after switching to the bytes solution. There is only one submitted solution at this point.';
    $wd.restoreSolution;
    $wd.isBytesAndChars: 210, 88, 'after discarding the bytes solution.';
    $wd.isSolutionPickerState: '', 'after discarding the bytes solution.';
}

subtest 'When not logged in, after a user submits different bytes and chars solutions and then logs in, they can submit the bytes solution and discard the chars solution.' => {
    plan 12;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Python').click;
    $wd.typeCode: $python210_88;
    $wd.isBytesAndChars: 210, 88, 'after typing the chars solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the chars solution, while logged out.';
    $wd.typeCode: BACKSPACE x 88 ~ $python121_121;
    $wd.isBytesAndChars: 121, 121, 'after typing the bytes solution.';
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged out.';
    $wd.isSolutionPickerState: 'bytes', 'after submitting the bytes solution, while logged out.';
    # Log in and reload the page.
    $wd.setSessionCookie: createUserAndSession;
    $wd.loadFizzBuzz;
    $wd.isBytesAndChars: 121, 121, 'after reloading the page';
    $wd.isSolutionPickerState: 'bytes', 'after reloading the page';
    # Submit the solution as a logged-in user.
    $wd.run;
    $wd.isPassing: 'after submitting the bytes solution, while logged in.';
    $wd.setSolution: 'chars';
    $wd.isBytesAndChars: 210, 88, 'after switching to the chars solution.';
    $wd.isSolutionPickerState: '', 'after switching to the chars solution. There is only one submitted solution at this point.';
    $wd.restoreSolution;
    $wd.isBytesAndChars: 121, 121, 'after discarding the chars solution.';
    $wd.isSolutionPickerState: '', 'after discarding the chars solution.';
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
    $wd.isBytesAndChars: 121, 121, 'after typing code.';
    $wd.run;
    $wd.isPassing: 'after running code.';
    $wd.isSolutionPickerState: '', 'after running code.';
    # Improve the solution outside of this browser session.
    ok post-solution(:code($python62_62), :hole<fizz-buzz>, :lang<python>, :$session)<Pass>, 'Passes';
    $wd.loadFizzBuzz;
    $wd.isRestoreSolutionLinkVisible: False, 'after improving the solution outside of the browser session and reloading.';
    $wd.isBytesAndChars: 62, 62, 'The byte count should be lower, after reloading the page.';
    $wd.isSolutionPickerState: '', 'after reloading the page.';
}
