FROM rust:1.93.1-slim-trixie AS builder

ENV CRT='x86_64-unknown-linux-musl' VER=0.18.0

RUN rustup target add $CRT

RUN cargo install uiua    \
    --features binary     \
    --no-default-features \
    --root /usr           \
    --target $CRT         \
    --version $VER        \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /usr/bin/uiua /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
