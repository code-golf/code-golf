FROM alpine:3.15 as builder

# Zig expects /dev/null to exist.
RUN mkdir /empty /mydev \
 && mknod -m 0666 /mydev/null c 1 3

RUN apk add --no-cache binutils curl

RUN curl https://ziglang.org/download/0.9.1/zig-linux-x86_64-0.9.1.tar.xz \
  | tar xJC tmp --strip 1 \
 && strip /tmp/zig

FROM scratch

COPY --from=0 /bin/cat /bin/rm /bin/sh /bin/
COPY --from=0 /mydev                   /dev
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /tmp/zig                 /usr/local/bin/
COPY --from=0 /tmp/lib                 /usr/local/lib

COPY zig /usr/bin/

ENTRYPOINT ["zig"]

CMD ["version"]
