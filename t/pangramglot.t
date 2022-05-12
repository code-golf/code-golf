use t;

for 'config/langs.toml'.IO.&from-toml.keys.sort: &lc -> $name {
    my $id = $name.lc.trans: qw[# + ><>] => qw[-sharp p fish];

    is dbh.execute(Q:s "SELECT pangramglot('{$id}')").row.head,
        ($name.lc.comb ∩ set('a' .. 'z')).elems, $name;
}

done-testing;
