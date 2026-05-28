FROM rust:1.95.0-slim-trixie AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes musl-tools patch

ENV CRT='x86_64-unknown-linux-musl' VER=0.113.0

WORKDIR /nushell

ADD https://github.com/nushell/nushell.git#$VER .

COPY nushell.patch .

RUN patch -p0 < nushell.patch

RUN rustup target add $CRT

RUN cargo install \
    --path .      \
    --root /usr   \
    --target $CRT \
 && strip /usr/bin/nu

FROM codegolf/lang-base

COPY --from=0 /usr/bin/nu /usr/bin/nushell

ENTRYPOINT ["nushell"]

CMD ["--version"]
