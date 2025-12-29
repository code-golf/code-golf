FROM rust:1.92.0-alpine3.23 AS builder

RUN apk add --no-cache curl musl-dev

ENV VER=b8d2274

RUN curl -#L https://github.com/betaveros/noulith/tarball/$VER \
  | tar xz --strip-components 1

RUN cargo install \
    --path /      \
    --root /usr   \
 && strip /usr/bin/noulith

FROM codegolf/lang-base

COPY --from=0 /usr/bin/noulith /usr/bin/

ENTRYPOINT ["noulith"]
