FROM alpine:3.19 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/tcltk/tcl/archive/refs/tags/core-8-6-13.tar.gz | tar xz

ENV LDFLAGS -static

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
