FROM alpine:3.16 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/tcltk/tcl/archive/refs/tags/core-8-6-12.tar.gz | tar xz

RUN touch /empty

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

# Create empty init.tcl to avoid "can't find usable init.tcl" error.
COPY --from=0 /empty          /library/init.tcl
COPY --from=0 /tcl/unix/tclsh /usr/bin/tcl

ENTRYPOINT ["tcl"]
