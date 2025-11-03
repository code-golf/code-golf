FROM rust:1.91.0-slim-trixie AS builder

ENV CRT='x86_64-unknown-linux-musl' VER=0.1.4

RUN rustup target add $CRT

RUN cargo install hush \
    --root /usr        \
    --target $CRT      \
    --version $VER     \
 && strip /usr/bin/hush

FROM codegolf/lang-base

COPY --from=0 /usr/bin/hush /usr/bin/

ENTRYPOINT ["hush"]

CMD ["--version"]
