FROM rust:1.87.0-slim-bookworm AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes build-essential

ENV RUSTFLAGS='-C target-feature=+crt-static' VER=0.16.0

RUN cargo install uiua                \
    --features binary                 \
    --no-default-features             \
    --root /usr                       \
    --target x86_64-unknown-linux-gnu \
    --version $VER                    \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /usr/bin/uiua /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
