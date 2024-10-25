FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base cmake curl ncurses-dev

ENV VERSION=b7bf640

RUN curl -L https://github.com/VorpalBlade/cfunge/tarball/$VERSION | tar xz \
 && cd /VorpalBlade-cfunge-$VERSION && mkdir bin && cd bin && cmake .. && make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh /bin/touch /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /*-cfunge-*/bin/cfunge      /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6   /usr/lib/

COPY /befunge /usr/bin/

ENTRYPOINT ["befunge"]

CMD ["--version"]
