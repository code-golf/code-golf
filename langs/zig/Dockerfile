FROM alpine:3.21 AS builder

RUN apk add --no-cache binutils curl gcc musl-dev

ENV VER=0.14.1

RUN curl -# https://ziglang.org/download/$VER/zig-x86_64-linux-$VER.tar.xz \
  | tar xJC tmp --strip-components 1

# Remove some platforms we don't need.
RUN rm -r /tmp/lib/libc/darwin \
          /tmp/lib/libc/glibc  \
          /tmp/lib/libc/mingw  \
          /tmp/lib/libc/wasi   \
 && strip /tmp/zig

COPY zig.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/zigwrapper -s -static zig.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/env             \
              /usr/bin/zigwrapper      /usr/bin/
COPY --from=0 /tmp/zig                 /usr/local/bin/
COPY --from=0 /tmp/lib                 /usr/local/lib

ENTRYPOINT ["zigwrapper"]

CMD ["--version"]
