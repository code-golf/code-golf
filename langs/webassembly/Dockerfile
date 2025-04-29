FROM rust:1.86.0-alpine3.21 AS builder

RUN apk add --no-cache musl-dev

ENV VER=31.0.0

RUN wget -O- https://github.com/bytecodealliance/wasmtime/archive/refs/tags/v$VER.tar.gz | tar xz

RUN cargo install         \
    --path /wasmtime-$VER \
    --root /usr           \
 && strip /usr/bin/wasmtime

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/wasmtime        /usr/bin/webassembly

ENTRYPOINT ["webassembly"]

CMD ["--version"]
