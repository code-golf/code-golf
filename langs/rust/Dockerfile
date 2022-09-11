FROM rust:1.63.0-alpine3.16 as builder

WORKDIR /usr/local/rustup/toolchains/1.63.0-x86_64-unknown-linux-musl

# Delete unused binaries.
RUN find /usr/local/cargo/bin bin -not -name rustc -delete \
 && rm lib/rustlib/x86_64-unknown-linux-musl/bin/rust-lld

# Delete unused libs.
RUN rm lib/libtest-*.so

# Strip the giant rust libs.
RUN strip lib/*.so

FROM codegolf/lang-base

COPY --from=0 /bin/sh                         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1        \
              /lib/libz.so.1                  /lib/
COPY --from=0 /usr/bin/ld                     /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.38.so         \
              /usr/lib/libctf.so.0            \
              /usr/lib/libgcc_s.so.1          /usr/lib/
COPY --from=0 /usr/local/cargo                /usr/local/cargo
COPY --from=0 /usr/local/rustup/settings.toml /usr/local/rustup/

COPY --from=0 /usr/local/rustup/toolchains/1.63.0-x86_64-unknown-linux-musl/bin \
              /usr/local/rustup/toolchains/1.63.0-x86_64-unknown-linux-musl/bin
COPY --from=0 /usr/local/rustup/toolchains/1.63.0-x86_64-unknown-linux-musl/lib \
              /usr/local/rustup/toolchains/1.63.0-x86_64-unknown-linux-musl/lib

COPY rust /usr/bin/

ENTRYPOINT ["rust"]

CMD ["--version"]
