# Also bump Gleam, Hush.
ARG VER=1.83.0

FROM rust:$VER-alpine3.21 AS builder

ARG VER

WORKDIR /usr/local/rustup/toolchains/$VER-x86_64-unknown-linux-musl

# Delete unused binaries.
RUN find /usr/local/cargo/bin bin -not -name rustc -delete \
 && rm lib/rustlib/x86_64-unknown-linux-musl/bin/rust-lld

# Delete unused sanitizers.
RUN find . -name '*.?san.a' -delete

# Strip the giant rust libs.
RUN strip lib/*.so

FROM codegolf/lang-base

ARG VER

COPY --from=0 /bin/sh                         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1        /lib/
COPY --from=0 /usr/bin/ld                     /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.43.1.so       \
              /usr/lib/libctf.so.0            \
              /usr/lib/libgcc_s.so.1          \
              /usr/lib/libjansson.so.4        \
              /usr/lib/libsframe.so.1         \
              /usr/lib/libz.so.1              \
              /usr/lib/libzstd.so.1           /usr/lib/
COPY --from=0 /usr/local/cargo                /usr/local/cargo
COPY --from=0 /usr/local/rustup/settings.toml /usr/local/rustup/

COPY --from=0 /usr/local/rustup/toolchains/$VER-x86_64-unknown-linux-musl/bin \
              /usr/local/rustup/toolchains/$VER-x86_64-unknown-linux-musl/bin
COPY --from=0 /usr/local/rustup/toolchains/$VER-x86_64-unknown-linux-musl/lib \
              /usr/local/rustup/toolchains/$VER-x86_64-unknown-linux-musl/lib

COPY rust /usr/bin/

ENTRYPOINT ["rust"]

CMD ["--version"]
