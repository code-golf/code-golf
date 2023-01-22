use t;

for 'config/langs.toml'.IO.&from-toml.sort {
    next if .value<experiment>;

    my $id = .key.lc.trans: qw[# + ><>] => qw[-sharp p fish];

    is dbh.execute(Q:s "SELECT pangramglot('{$id}')").row.head,
        (.key.lc.comb ∩ set('a' .. 'z')).elems, .key;
}

done-testing;
