FROM rust:1.91.1-slim-trixie AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes bash curl patch

ENV CRT='x86_64-unknown-linux-musl' VER=0bcafce

WORKDIR /hurl

RUN curl -#L https://git.sr.ht/~ntietz/hurl-lang/archive/$VER.tar.gz \
  | tar xz --strip-components 1

COPY hurl.patch .

RUN rustup target add $CRT \
 && sed -i '1d' src/lib.rs \
 && patch -p0 < hurl.patch

RUN cargo install     \
    --path .          \
    --root /usr/local \
    --target $CRT     \
 && strip /usr/local/bin/hurl

COPY hurl.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/hurl -s -static /hurl.c

FROM codegolf/lang-base

COPY --from=0 /hurl/lib           /
COPY --from=0 /usr/bin/hurl       /usr/bin/
COPY --from=0 /usr/local/bin/hurl /usr/local/bin/
