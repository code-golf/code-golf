FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=3.9.7

RUN curl -#L https://downloads.sourceforge.net/regina-rexx/regina-rexx-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure      \
    --disable-flists \
    --disable-gci    \
    --disable-screen \
    --enable-64bit   \
    --prefix=/usr    \
 && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /usr/bin/regina         /usr/bin/rexx
COPY --from=0 /usr/lib/libregina.so.3 /usr/lib/

ENTRYPOINT ["rexx"]

CMD ["--version"]
