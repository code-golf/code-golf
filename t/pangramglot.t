use t;

for 'config/data/langs.toml'.IO.&from-toml.sort {
    my $id = .key.lc.subst(' ', '-').trans: qw[# + ><>] => qw[-sharp p fish];

    is dbh.execute(Q:s "SELECT pangramglot('{$id}')").row.head,
        (.key.lc.comb ∩ set('a' .. 'z')).elems, .key;
}

done-testing;
