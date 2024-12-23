FROM rust:1.83.0-alpine3.21 AS builder

RUN apk add --no-cache build-base linux-headers

RUN cargo install uiua --version 0.14.0 \
 && strip /usr/local/cargo/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/local/cargo/bin/uiua /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
