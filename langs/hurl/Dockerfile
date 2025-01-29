ARG VER=0bcafce7

FROM rust:1.84.0-alpine3.21 AS builder

ARG VER

RUN apk add --no-cache build-base curl

RUN curl -# https://git.sr.ht/~ntietz/hurl-lang/archive/$VER.tar.gz | tar xz

WORKDIR /hurl-lang-$VER

COPY lib_parser.patch /

RUN patch -p1 < /lib_parser.patch

RUN cargo install --path . \
 && strip /usr/local/cargo/bin/hurl

COPY hurl.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/hurl -s /hurl.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /hurl-lang-$VER/lib       /
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/hurl             /usr/bin/
COPY --from=0 /usr/local/cargo/bin/hurl /usr/local/bin/

ENTRYPOINT ["hurl"]

CMD ["--version"]
