FROM rust:1.87.0-alpine3.22 AS builder

RUN apk add --no-cache build-base linux-headers

ENV VER=0.16.0

RUN cargo install uiua    \
    --features binary     \
    --no-default-features \
    --root /usr           \
    --version $VER        \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/uiua            /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
