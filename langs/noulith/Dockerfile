FROM rust:1.91.1-alpine3.22 AS builder

RUN apk add --no-cache curl musl-dev

ENV VER=b5703eb

RUN curl -#L https://github.com/betaveros/noulith/tarball/$VER \
  | tar xz --strip-components 1

RUN cargo install \
    --path /      \
    --root /usr   \
 && strip /usr/bin/noulith

FROM codegolf/lang-base

COPY --from=0 /usr/bin/noulith /usr/bin/

ENTRYPOINT ["noulith"]
