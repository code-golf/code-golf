FROM rust:1.82.0-alpine3.20 AS builder

RUN apk add --no-cache build-base

RUN cargo install hush --version 0.1.4

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/local/cargo/bin/hush /usr/bin/

ENTRYPOINT ["hush"]

CMD ["--version"]
