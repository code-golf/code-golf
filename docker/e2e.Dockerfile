FROM rakuland/raku

RUN apk add --no-cache libpq openssl-dev \
 && zef --/test install                  \
    'DBIish:ver<0.6.3>'                  \
    'IO::Socket::SSL:ver<0.0.2>'         \
    'JSON::Fast:ver<0.16>'               \
    'TOML::Thumb:ver<0.2>'               \
    git://gitlab.com/JRaspass/webdriver.git@0e29ee39
