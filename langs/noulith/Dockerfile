FROM rust:1.80.1-alpine3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/betaveros/noulith/archive/828de52.tar.gz \
  | tar xz --strip-components 1

RUN cargo build    \
    --features cli \
    --release      \
 && cargo fix      \
    --allow-no-vcs \
    --lib          \
    --package noulith

FROM codegolf/lang-base

COPY --from=0 /target/release/noulith /usr/bin/

ENTRYPOINT ["noulith"]
