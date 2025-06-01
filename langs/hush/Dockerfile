FROM rust:1.87.0-alpine3.22 AS builder

RUN apk add --no-cache musl-dev

ENV VER=0.1.4

RUN cargo install hush \
    --root /usr        \
    --version $VER

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/hush            /usr/bin/

ENTRYPOINT ["hush"]

CMD ["--version"]
