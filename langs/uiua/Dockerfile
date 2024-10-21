FROM rust:1.82.0-alpine3.20 AS builder

RUN apk add --no-cache build-base linux-headers

RUN cargo install uiua --version 0.12.3

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/local/cargo/bin/uiua /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
