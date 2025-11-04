FROM rust:1.91.0-slim-trixie AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl musl-dev

ENV CRT='x86_64-unknown-linux-musl' VER=38.0.3

WORKDIR /usr/bin

RUN curl -#L https://github.com/bytecodealliance/wasmtime/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN rustup target add $CRT

RUN cargo install \
    --path .      \
    --root /usr   \
    --target $CRT \
 && strip wasmtime

FROM codegolf/lang-base

COPY --from=0 /usr/bin/wasmtime /usr/bin/webassembly

ENTRYPOINT ["webassembly"]

CMD ["--version"]
