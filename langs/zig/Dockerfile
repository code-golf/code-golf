FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=0.15.1

RUN curl -#L https://ziglang.org/download/$VER/zig-x86_64-linux-$VER.tar.xz \
  | tar xJ --directory /usr/local --strip-components 1

COPY zig.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/zig -s -static zig.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1       /lib/
COPY --from=0 /usr/bin/env /usr/bin/zig      /usr/bin/
COPY --from=0 /usr/local/zig                 /usr/local/
COPY --from=0 /usr/local/lib/compiler_rt.zig \
              /usr/local/lib/ubsan_rt.zig    /usr/local/lib/
COPY --from=0 /usr/local/lib/compiler_rt     /usr/local/lib/compiler_rt
COPY --from=0 /usr/local/lib/std             /usr/local/lib/std

ENTRYPOINT ["zig"]

CMD ["--version"]
