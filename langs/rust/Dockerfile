FROM rust:1.62.1-alpine3.16 as builder

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.62.1-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

# Strip the giant rust libs.
RUN strip /usr/local/rustup/toolchains/1.62.1-x86_64-unknown-linux-musl/lib/*.so

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

COPY --from=0 /usr/local/rustup/toolchains/1.62.1-x86_64-unknown-linux-musl/bin \
              /usr/local/rustup/toolchains/1.62.1-x86_64-unknown-linux-musl/bin
COPY --from=0 /usr/local/rustup/toolchains/1.62.1-x86_64-unknown-linux-musl/lib \
              /usr/local/rustup/toolchains/1.62.1-x86_64-unknown-linux-musl/lib

COPY rust /usr/bin/

ENTRYPOINT ["rust"]

CMD ["--version"]
