FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=5.3.2

RUN curl -#L https://ftpmirror.gnu.org/gnu/gawk/gawk-$VER.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure    \
    --disable-lint \
    --prefix=/usr  \
 && make install   \
 && strip /usr/bin/gawk

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/gawk       /usr/bin/awk
COPY --from=0 /usr/lib/gawk       /usr/lib/gawk

ENTRYPOINT ["awk"]

CMD ["--version"]
