FROM rust:1.56-alpine3.13 as builder

RUN mkdir /empty

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.56.0-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

# Delete some large unneeded libs.
RUN rm                             \
    /lib/libcrypto.so.1.1          \
    /lib/libssl.so.1.1             \
    /usr/lib/libgdruntime.so.1.0.0 \
    /usr/lib/libgphobos.so.1.0.0   \
    /usr/lib/libisl.so.22.0.0      \
    /usr/lib/libmpfr.so.6.1.0      \
    /usr/lib/libopcodes-2.35.2.so  \
    /usr/lib/libstdc++.so.6.0.28

# Strip the giant rust libs.
RUN strip /usr/local/rustup/toolchains/1.56.0-x86_64-unknown-linux-musl/lib/*.so

FROM scratch

ENV RUSTUP_HOME /usr/local/rustup

# /bin/cc (gcc) seems to be needed as a linker :-(
COPY --from=0 /bin/sh                       /bin/
COPY --from=0 /lib                          /lib
COPY --from=0 /empty                        /proc
COPY --from=0 /empty                        /tmp
COPY --from=0 /usr/bin/cc                   /usr/bin/
COPY --from=0 /usr/lib                      /usr/lib
COPY --from=0 /usr/libexec                  /usr/libexec
COPY --from=0 /usr/local/cargo              /usr/local/cargo
COPY --from=0 /usr/local/rustup             /usr/local/rustup
COPY --from=0 /usr/x86_64-alpine-linux-musl /usr/x86_64-alpine-linux-musl

COPY rust /usr/bin/

ENTRYPOINT ["/usr/bin/rust"]

CMD ["--version"]
