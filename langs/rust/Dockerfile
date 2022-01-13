FROM rust:1.58.0-alpine3.15 as builder

RUN mkdir /empty

# Delete the non rustc binaries.
RUN find /usr/local/cargo/bin \
    /usr/local/rustup/toolchains/1.58.0-x86_64-unknown-linux-musl/bin \
    -not -name rustc -delete

# Strip the giant rust libs.
RUN strip /usr/local/rustup/toolchains/1.58.0-x86_64-unknown-linux-musl/lib/*.so

FROM scratch

ENV RUSTUP_HOME /usr/local/rustup/10.3.1

COPY --from=0 /bin/sh                       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1      \
              /lib/libz.so.1                /lib/
COPY --from=0 /empty                        /proc
COPY --from=0 /empty                        /tmp
COPY --from=0 /usr/bin/ld                   /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.37.so       \
              /usr/lib/libctf.so.0          \
              /usr/lib/libgcc_s.so.1        /usr/lib/
COPY --from=0 /usr/local/cargo              /usr/local/cargo
COPY --from=0 /usr/local/rustup             /usr/local/rustup
COPY --from=0 /usr/x86_64-alpine-linux-musl /usr/x86_64-alpine-linux-musl

COPY rust /usr/bin/

ENTRYPOINT ["rust"]

CMD ["--version"]
