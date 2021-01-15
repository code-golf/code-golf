FROM alpine:3.13 as builder

RUN mkdir /empty

RUN apk add --no-cache binutils curl

RUN curl https://ziglang.org/download/0.7.0/zig-linux-x86_64-0.7.0.tar.xz \
  | tar -xJf -

RUN strip /zig-linux-x86_64-0.7.0/zig

FROM scratch

COPY --from=0 /bin/cat /bin/rm /bin/sh    /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /empty                      /proc
COPY --from=0 /empty                      /tmp
COPY --from=0 /zig-linux-x86_64-0.7.0/zig /usr/local/bin/
COPY --from=0 /zig-linux-x86_64-0.7.0/lib /usr/local/lib/

COPY zig /usr/bin/

ENTRYPOINT ["/usr/local/bin/zig", "version"]
