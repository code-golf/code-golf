FROM rust:1.41.1-alpine3.11 as builder

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.41.1-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

RUN rm -r                        \
    /lib/apk                     \
    /lib/firmware                \
    /lib/libcrypto.*             \
    /lib/libssl.*                \
    /lib/mdev                    \
    /usr/bin/gcov*               \
    /usr/bin/ssl_client          \
    /usr/include                 \
    /usr/lib/libcrypto.*         \
    /usr/lib/libssl.*            \
    /usr/lib/libtls-standalone.* \
    /usr/local/bin               \
    /usr/local/lib               \
    /usr/local/rustup/downloads  \
    /usr/local/rustup/tmp        \
    /usr/local/share             \
    /usr/sbin                    \
    /usr/share

FROM scratch

ENV RUSTUP_HOME /usr/local/rustup

COPY --from=0 /bin/sh /bin/
COPY --from=0 /lib    /lib/
COPY --from=0 /usr    /usr/

COPY /compile_and_run.sh /usr/bin/rust

ENTRYPOINT [ "/usr/local/cargo/bin/rustc", "--version" ]
