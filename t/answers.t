use t;

for sort dir 'hole/answers' {
    my $code = 'say ｢' ~ .slurp ~ '｣';
    my $hole = .extension("").basename;

    $hole = '√2' if $hole eq 'root-2';

    ok post-solution( :$code :$hole )<Pass>, $hole;
}

done-testing;
