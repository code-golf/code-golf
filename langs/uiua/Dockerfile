FROM rust:1.87.0-slim-bookworm AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes gcc make

ENV VER=0.16.2

RUN RUSTFLAGS='-C target-feature=+crt-static' \
    cargo install uiua                        \
    --features binary                         \
    --git https://github.com/uiua-lang/uiua   \
    --no-default-features                     \
    --root /usr                               \
    --tag $VER                                \
    --target x86_64-unknown-linux-gnu         \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /usr/bin/uiua /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
