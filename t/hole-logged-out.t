# This file contains tests for the hole route where not being logged in is an important characteristic.
# Some of these tests have similar counterparts in hole-logged-in.t.

use t;
use hole;

plan 8;

subtest 'Successful solutions are loaded from localStorage on reload.' => {
    plan 5;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.typeCode: $raku57_55;
    $wd.isBytesAndChars: 57, 55;
    $wd.run;
    $wd.isPassing;
    $wd.isSolutionPickerState: '';
    $wd.clearLocalStorage;
    $wd.loadFizzBuzz;
    $wd.getLangLink('Raku').click;
    $wd.isBytesAndChars: 0, 0, 'after clearing localStorage and reloading the page.';
    $wd.isSolutionPickerState: '', 'after clearing localStorage and reloading the page.';
}

subtest 'Untested solutions are loaded from localStorage on reload.' => {
    plan 5;
    my $wd = HoleWebDriver.create;
    LEAVE $wd.delete-session;
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

subtest 'Different bytes and chars solutions, and the active solution, are loaded from localStorage on reload.' => {
    plan 11;
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
}
