FROM alpine:3.24 AS builder

ENV VER=0.17.1

ADD --unpack https://ziglang.org/download/$VER/zig-x86_64-linux-$VER.tar.xz /

FROM codegolf/lang-base

COPY --from=0           /lib/ld-musl-*.so.1          /lib/
COPY --from=0           /zig-*/zig /usr/bin/env      /usr/bin/
COPY --from=0 --parents /zig-*/lib/./compiler_rt     \
                        /zig-*/lib/./compiler_rt.zig \
                        /zig-*/lib/./std             \
                        /zig-*/lib/./ubsan_rt.zig    /usr/lib/zig

ENTRYPOINT ["zig"]

CMD ["version"]
