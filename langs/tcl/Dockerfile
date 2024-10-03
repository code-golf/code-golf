FROM alpine:3.20 AS builder

RUN apk add --no-cache bsd-compat-headers build-base curl zlib-dev

RUN curl -L https://github.com/tcltk/tcl/archive/refs/tags/core-9-0-0.tar.gz | tar xz

ENV LDFLAGS=-static

RUN mv /tcl-* /tcl   \
 && cd tcl/unix      \
 && ./configure      \
    --disable-load   \
    --disable-shared \
    --enable-64-bit  \
 && make             \
 && strip tclsh

FROM codegolf/lang-base

COPY --from=0 /tcl/library/   /library/
COPY --from=0 /tcl/unix/tclsh /usr/bin/tcl

ENTRYPOINT ["tcl"]
