FROM alpine:3.16 as builder

RUN apk add --no-cache binutils curl

RUN curl https://ziglang.org/download/0.10.0/zig-linux-x86_64-0.10.0.tar.xz \
  | tar xJC tmp --strip 1 \
 && strip /tmp/zig

# Remove some platforms we don't need.
RUN rm -r /tmp/lib/libc/darwin \
          /tmp/lib/libc/glibc  \
          /tmp/lib/libc/mingw  \
          /tmp/lib/libc/wasi

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/rm /bin/sh /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/env             /usr/bin/
COPY --from=0 /tmp/zig                 /usr/local/bin/
COPY --from=0 /tmp/lib                 /usr/local/lib

COPY zig /usr/bin/

ENTRYPOINT ["/usr/bin/zig"]

CMD ["version"]
