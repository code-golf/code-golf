FROM rust:1.83.0-alpine3.21 AS builder

RUN apk add --no-cache musl-dev

ENV VER=57580ee

RUN wget -O- https://github.com/betaveros/noulith/archive/$VER.tar.gz \
  | tar xz --strip-components 1

RUN cargo build    \
    --features cli \
    --release      \
 && cargo fix      \
    --allow-no-vcs \
    --lib          \
    --package noulith

RUN strip /target/release/noulith

FROM codegolf/lang-base

COPY --from=0 /target/release/noulith /usr/bin/

ENTRYPOINT ["noulith"]
