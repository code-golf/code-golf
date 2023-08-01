FROM alpine:3.18 AS builder

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing --no-cache \
    curl gcc musl-dev pcre tcc

ENV VER 2.0.0

RUN curl https://nim-lang.org/download/nim-$VER.tar.xz | tar xJ

RUN cd nim-$VER             \
 && ./build.sh              \
 && ./bin/nim c koch        \
 && ./koch boot -d:release  \
 && ./install.sh /          \
 && strip /nim/bin/nim

# Docker can only copy symlinks by copying the whole folder :-(
RUN mkdir /usr-lib-symlinks && cp -P /usr/lib/libc.so /usr-lib-symlinks

FROM codegolf/lang-base

COPY --from=0 /bin/sh                  /bin/
COPY          /nim.cfg                 /etc/nim/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /nim/bin/nim             \
              /usr/bin/tcc             /usr/bin/
COPY --from=0 /usr/include             /usr/include
COPY --from=0 /usr-lib-symlinks        /usr/lib
COPY --from=0 /usr/lib/crt1.o          \
              /usr/lib/crti.o          \
              /usr/lib/crtn.o          \
              /usr/lib/libdl.a         \
              /usr/lib/libm.a          \
              /usr/lib/libpcre.so.1    \
              /usr/lib/librt.a         /usr/lib/
COPY --from=0 /nim/lib                 /usr/lib/nim/lib
COPY --from=0 /usr/lib/tcc             /usr/lib/tcc

ENTRYPOINT ["nim"]

CMD ["-v"]
