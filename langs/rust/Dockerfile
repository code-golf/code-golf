FROM rust:1.41.1-alpine3.11 as builder

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.41.1-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

FROM scratch

ENV RUSTUP_HOME /usr/local/rustup

COPY --from=0 / /

COPY /compile_and_run.sh /usr/bin/rust

ENTRYPOINT [ "/usr/local/cargo/bin/rustc", "--version" ]
