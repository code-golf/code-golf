use t;

for (

    # System.
    # Use Raku because our Perl lacks the "Sys::Hostname" package.
    cwd      => '/',              raku => 'say ~$*CWD',
    env      => '{HOME => /tmp}', raku => 'say %*ENV',
    hostname => 'code-golf',      raku => 'say $*KERNEL.hostname',

    # User & Group.
    # Use Perl because our Raku lacks the "id" binary.
    group    => 'nobody', perl => 'say +( getgrgid $( )[0]',
    user     => 'nobody', perl => 'say +( getpwuid $< )[0]',
    group-id => '65534',  perl => 'say 0 + $(',
    user-id  => '65534',  perl => 'say 0 + $<',

) -> ( :key($name), :value($exp) ), ( :key($lang), :value($code) ) {
    my $got = post-solution( :$code :$lang )<runs>[0];

    is $got<stdout>, $exp, $name;

    diag $got<stderr> if $got<stderr>;
}

done-testing;
