FROM alpine:3.22 AS builder

RUN apk add --no-cache bsd-compat-headers build-base curl zlib-dev

ENV LDFLAGS='-static' VER=9-0-1

RUN curl -#L https://github.com/tcltk/tcl/archive/refs/tags/core-$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /unix

RUN ./configure      \
    --disable-load   \
    --disable-shared \
    --enable-64-bit  \
    --prefix=/usr    \
 && make install     \
 && strip /usr/bin/tcl*

FROM codegolf/lang-base

COPY --from=0 /library      /library
COPY --from=0 /usr/bin/tcl* /usr/bin/tcl

ENTRYPOINT ["tcl"]

CMD ["--version"]
