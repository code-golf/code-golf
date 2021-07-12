FROM rakuland/raku

RUN apk add --no-cache libpq openssl-dev   \
 && zef --/test install DBIish             \
    IO::Socket::SSL JSON::Fast TOML::Thumb \
    git://gitlab.com/JRaspass/webdriver.git@0e29ee39
