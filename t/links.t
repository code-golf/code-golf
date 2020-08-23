use t;

=begin TODO
for (
    sort map $_->{url} =~ s|^/|https:/|r, map @{ $_->{links} // [] },
    values from_toml( read_binary 'holes.toml' )->%*
) {
    $_ = $ua->head($_);

    is $_->{status}, 301, "301 $_->{url}" for $_->{redirects}->@*;
    is $_->{status}, 200, "200 $_->{url}";
}
=end TODO

done-testing;
