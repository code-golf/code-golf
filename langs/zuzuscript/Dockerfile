FROM rust:1.96.0-slim-trixie AS builder

RUN apt-get update \
 && apt-get install --yes libclang-dev libffi-dev libssl-dev libxml2-dev pkgconf

ENV VER=0.6.0

ADD git@github.com:zuzuscript/zuzu-rust.git#v$VER /zuzu

WORKDIR /zuzu

RUN cargo install --path . --root /usr \
 && ./post-install.sh --prefix /usr    \
 && strip /usr/bin/zuzu-rust

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libcrypto.so.3 \
              /lib/x86_64-linux-gnu/libffi.so.8    \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/liblzma.so.5   \
              /lib/x86_64-linux-gnu/libm.so.6      \
              /lib/x86_64-linux-gnu/libxml2.so.2   \
              /lib/x86_64-linux-gnu/libssl.so.3    \
              /lib/x86_64-linux-gnu/libz.so.1      \
              /lib/x86_64-linux-gnu/libzstd.so.1   /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=0 /usr/bin/zuzu-rust                   /usr/bin/zuzu
COPY --from=0 /usr/share/zuzu-rust/modules         /usr/share/zuzu-rust/modules

ENTRYPOINT ["zuzu"]

CMD ["-v"]
