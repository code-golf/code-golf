FROM alpine:3.13 as builder

ENV VERSION=f749e59

RUN apk add --no-cache curl

RUN curl -L https://github.com/m-ender/hexagony/tarball/$VERSION | tar xz \
 && mv /m-ender-hexagony-$VERSION /hexagony

FROM codegolf/lang-ruby

COPY --from=0 /hexagony /usr/bin/hexagony

COPY interpreter.rb /usr/bin/hexagony/
