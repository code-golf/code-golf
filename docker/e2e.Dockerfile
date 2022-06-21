FROM rakuland/raku

RUN apk add --no-cache libpq openssl-dev && zef --/test install \
    'DBIish:ver<0.6.4>:auth<zef:raku-community-modules>'        \
    'HTTP::Tiny:ver<0.1.10>:auth<zef:jjatria>'                  \
    'IO::Socket::SSL:ver<0.0.2>:auth<github:sergot>'            \
    'JSON::Fast:ver<0.17>:auth<cpan:TIMOTIMO>'                  \
    'TOML::Thumb:ver<0.2>:auth<zef:JRaspass>'                   \
    'WebDriver:ver<0.1>:auth<zef:raku-land>'
