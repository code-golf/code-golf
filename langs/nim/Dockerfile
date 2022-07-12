FROM alpine:3.16 as builder

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing --no-cache \
    curl gcc musl-dev pcre tcc

RUN curl https://nim-lang.org/download/nim-1.6.6.tar.xz | tar xJ

RUN cd nim-1.6.6            \
 && ./build.sh              \
 && ./bin/nim c koch        \
 && ./koch boot -d:release  \
 && ./koch install /usr/bin \
 && strip /usr/bin/nim

# Docker can only copy symlinks by copying the whole folder :-(
RUN mkdir /usr-lib-symlinks && cp -P /usr/lib/libc.so /usr-lib-symlinks

FROM codegolf/lang-base

COPY --from=0 /bin/sh                  /bin/
COPY          /nim.cfg                 /etc/nim/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/nim             \
              /usr/bin/tcc             /usr/bin/
COPY --from=0 /usr/include             /usr/include/
COPY --from=0 /usr-lib-symlinks        /usr/lib
COPY --from=0 /usr/lib/crt1.o          \
              /usr/lib/crti.o          \
              /usr/lib/crtn.o          \
              /usr/lib/libdl.a         \
              /usr/lib/libm.a          \
              /usr/lib/libpcre.so.1    \
              /usr/lib/librt.a         /usr/lib/
COPY --from=0 /usr/lib/nim             /usr/lib/nim/
COPY --from=0 /usr/lib/tcc             /usr/lib/tcc/

ENTRYPOINT ["nim"]

CMD ["-v"]
