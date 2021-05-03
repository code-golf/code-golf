use t;

constant $pids24 = '
    say $$;
    for (1..23) { next if fork // die; say $$; sleep 1; exit }
    1 while wait != -1
';

constant $pids25 = '
    say $$;
    for (1..24) { next if fork // die; say $$; sleep 1; exit }
    1 while wait != -1
';

subtest '24 PIDs succeed' => {
    my $res = post-solution code => $pids24, lang => 'perl';

    is $res<ExitCode>, 0, 'ExitCode';
    is $res<Out>.lines.sort(+*), 1..24, 'Out';
};

subtest '25 PIDs fail' => {
    my $res = post-solution code => $pids25, lang => 'perl';

    is $res<ExitCode>, 11, 'ExitCode';
    isnt $res<Out>.lines.sort(+*), 1..25, 'Out';
};

done-testing;
