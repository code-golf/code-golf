FROM rust:1.47.0-alpine3.12 as builder

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.47.0-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

FROM scratch

ENV RUSTUP_HOME /usr/local/rustup

COPY --from=0 /     /
COPY          /rust /usr/bin

ENTRYPOINT [ "/usr/local/cargo/bin/rustc", "--version" ]
