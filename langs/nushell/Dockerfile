FROM rust:1.96.0-slim-trixie AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes musl-dev

ENV CRT='x86_64-unknown-linux-musl' VER=0.113.1

RUN rustup target add $CRT

RUN cargo install nu      \
    --no-default-features \
    --root /usr           \
    --target $CRT         \
    --version $VER        \
 && strip /usr/bin/nu

FROM codegolf/lang-base

COPY --from=0 /usr/bin/nu /usr/bin/nushell

ENTRYPOINT ["nushell"]

CMD ["--version"]
