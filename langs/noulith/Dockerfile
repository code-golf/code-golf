FROM rust:1.84.0-alpine3.21 AS builder

RUN apk add --no-cache musl-dev

ENV VER=57580ee

RUN wget -O- https://github.com/betaveros/noulith/tarball/$VER \
  | tar xz --strip-components 1

RUN cargo install --path / --root /usr \
 && strip /usr/bin/noulith

FROM codegolf/lang-base

COPY --from=0 /usr/bin/noulith /usr/bin/

ENTRYPOINT ["noulith"]
