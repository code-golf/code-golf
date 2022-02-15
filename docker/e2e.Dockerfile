FROM rakuland/raku

RUN apk add --no-cache libpq openssl-dev && zef --/test install \
    'DBIish:ver<0.6.3>:auth<github:raku-community-modules>'     \
    'IO::Socket::SSL:ver<0.0.2>:auth<github:sergot>'            \
    'JSON::Fast:ver<0.17>:auth<cpan:TIMOTIMO>'                  \
    'TOML::Thumb:ver<0.2>:auth<zef:JRaspass>'                   \
    'WebDriver:ver<0.1>:auth<zef:raku-land>'
