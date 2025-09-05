# Also bump Gleam, Hush, Uiua.
ARG VER=1.89.0

FROM rust:$VER-alpine3.22 AS builder

ARG VER

RUN apk add --no-cache musl-dev

WORKDIR /usr/local/rustup/toolchains/$VER-x86_64-unknown-linux-musl

# Delete unused binary.
RUN rm lib/rustlib/x86_64-unknown-linux-musl/bin/rust-lld

# Delete unused sanitizers.
RUN find . -name '*.?san.a' -delete

# Strip the giant rust libs.
RUN strip lib/*.so

COPY rust.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rust -s -static /rust.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1                 /lib/
COPY --from=0 /usr/bin/ld /usr/bin/rust                /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.44.so                  \
              /usr/lib/libctf.so.0                     \
              /usr/lib/libgcc_s.so.1                   \
              /usr/lib/libjansson.so.4                 \
              /usr/lib/libsframe.so.1                  \
              /usr/lib/libz.so.1                       \
              /usr/lib/libzstd.so.1                    /usr/lib/
COPY --from=0 /usr/local/rustup/toolchains/*/bin/rustc /usr/local/bin/
COPY --from=0 /usr/local/rustup/toolchains/*/lib       /usr/local/lib

ENTRYPOINT ["rust"]

CMD ["--version"]
