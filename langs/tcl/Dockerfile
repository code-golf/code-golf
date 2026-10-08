FROM alpine:3.24 AS builder

RUN apk add --no-cache bsd-compat-headers gcc make musl-dev zlib-dev

ENV LDFLAGS='-static' VER=9.1.0

ADD --unpack https://sourceforge.net/projects/tcl/files/Tcl/$VER/tcl-core$VER-src.tar.gz/download /

WORKDIR /tcl$VER/unix

RUN ./configure      \
    --disable-load   \
    --disable-shared \
    --enable-64bit   \
    --prefix=/usr    \
 && make install     \
 && strip /usr/bin/tcl*

FROM codegolf/lang-base

COPY --from=0 /tcl*/library /usr/lib/tcl9.1
COPY --from=0 /usr/bin/tcl* /usr/bin/tcl

ENTRYPOINT ["tcl"]

CMD ["--version"]
