FROM alpine:3.16 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/tcltk/tcl/archive/refs/tags/core-8-6-12.tar.gz | tar xz

RUN touch /empty

RUN cd tcl-core-8-6-12/unix \
    && ./configure          \
    --disable-load          \
    --disable-symbols       \
    --disable-shared        \
    --enable-64-bit         \
    && make                 \
    && strip tclsh

FROM codegolf/lang-base

COPY --from=builder /bin/cat /bin/echo /bin/sh /bin/
COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder tcl-core-8-6-12/unix/tclsh /usr/bin/

# create empty init.tcl to avoid "can't find usable init.tcl" error
COPY --from=builder /empty /library/init.tcl

COPY tcl /usr/bin/

ENTRYPOINT ["tcl"]

CMD ["--version"]
