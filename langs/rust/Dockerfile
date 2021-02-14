FROM rust:1.50.0-alpine3.13 as builder

RUN mkdir /empty

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.50.0-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

FROM scratch

ENV RUSTUP_HOME /usr/local/rustup

COPY --from=0 /bin   /bin
COPY --from=0 /lib   /lib
COPY --from=0 /usr   /usr
COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY          /rust  /usr/bin

ENTRYPOINT [ "/usr/local/cargo/bin/rustc", "--version" ]
