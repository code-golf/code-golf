use t;

new-golfer :dbh(my $dbh = dbh);

is save-solution(:code<qwerty>, :time(123)), 123,
    'Initial solution';

is save-solution(:code<qwerty>, :time(456)), 123,
    "Slower identical solution didn't update runtime";

is save-solution(:code<qwerty>, :time(101)), 101,
    'Faster identical solution did update runtime';

is save-solution(:code<cool>, :time(789)), 789,
    'Slower shorter solution did update runtime';

todo 'FIXME';
is save-solution(:code<bean>, :time(999)), 999,
    "Slower same length solution did update runtime";

done-testing;

sub save-solution(:$code, :$time) {
    $dbh.execute:
        "SELECT save_solution(
            bytes := ?, chars := ?, code := ?, hole := 'fizz-buzz',
            lang  := 'perl', time_ms := ?::smallint, user_id := 1)",
        $code.encode('UTF-8').bytes, $code.chars, $code, $time;

    return $dbh.execute('SELECT time_ms FROM solutions LIMIT 1').row;
}
