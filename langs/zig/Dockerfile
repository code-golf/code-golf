FROM alpine:3.13 as builder

# Zig expects /dev/null to exist.
RUN mkdir /empty /mydev \
 && mknod -m 0666 /mydev/null c 1 3

RUN apk add --no-cache binutils curl

RUN curl https://ziglang.org/download/0.8.0/zig-linux-x86_64-0.8.0.tar.xz \
  | tar -xJ

RUN strip /zig-linux-x86_64-0.8.0/zig

FROM scratch

COPY --from=0 /bin/cat /bin/rm /bin/sh    /bin/
COPY --from=0 /mydev                      /dev
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /empty                      /proc
COPY --from=0 /empty                      /tmp
COPY --from=0 /zig-linux-x86_64-0.8.0/zig /usr/local/bin/
COPY --from=0 /zig-linux-x86_64-0.8.0/lib /usr/local/lib/

COPY zig /usr/bin/

ENTRYPOINT ["/usr/local/bin/zig", "version"]
