FROM rust:1.88.0-slim-bookworm AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl

ENV CRT='x86_64-unknown-linux-musl' VER=0.16.2

RUN curl -#L https://github.com/uiua-lang/uiua/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN rustup target add $CRT \
 && cargo add ffi

RUN cargo install         \
    --features binary     \
    --no-default-features \
    --path /              \
    --root /usr           \
    --target $CRT         \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /usr/bin/uiua /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
