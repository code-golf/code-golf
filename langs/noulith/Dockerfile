FROM rust:1.82.0-alpine3.20 AS builder

RUN apk add --no-cache build-base curl

ENV VERSION=57580ee

RUN curl -L https://github.com/betaveros/noulith/archive/$VERSION.tar.gz \
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
