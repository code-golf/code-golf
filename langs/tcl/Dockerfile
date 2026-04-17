FROM alpine:3.23 AS builder

RUN apk add --no-cache bsd-compat-headers curl gcc make musl-dev zlib-dev

ENV LDFLAGS='-static' VER=9.0.3

RUN curl -#L https://downloads.sourceforge.net/tcl/tcl-core$VER-src.tar.gz \
  | tar xz --strip-components 1

WORKDIR /unix

RUN ./configure      \
    --disable-load   \
    --disable-shared \
    --enable-64bit   \
    --prefix=/usr    \
 && make install     \
 && strip /usr/bin/tcl*

FROM codegolf/lang-base

COPY --from=0 /library      /library
COPY --from=0 /usr/bin/tcl* /usr/bin/tcl

ENTRYPOINT ["tcl"]

CMD ["--version"]
