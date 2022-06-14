FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl openssl-dev

ENV VERSION=802e45e

RUN curl -L https://github.com/lynn/golfscript/tarball/$VERSION | tar xz \
 && cd /lynn-golfscript-$VERSION && mv golfscript.rb /

FROM codegolf/lang-ruby

COPY --from=0 /golfscript.rb /usr/bin/golfscript

ENTRYPOINT ["golfscript"]

