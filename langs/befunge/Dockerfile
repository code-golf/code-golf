FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base cmake curl ncurses-dev

ENV VER=b7bf640

RUN curl -#L https://github.com/VorpalBlade/cfunge/tarball/$VER \
  | tar xz --strip-components 1

COPY strncpy.patch /

RUN patch -p 0 < strncpy.patch

WORKDIR /build

RUN cmake .. && make install \
 && strip /usr/local/bin/cfunge

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh /bin/touch /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /usr/local/bin/cfunge       /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6   /usr/lib/

COPY befunge /usr/bin/

ENTRYPOINT ["befunge"]

CMD ["--version"]
