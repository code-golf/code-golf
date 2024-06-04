FROM alpine:3.20 AS builder

RUN apk add --no-cache curl

RUN curl -L https://github.com/gleam-lang/gleam/releases/download/v1.2.1/gleam-v1.2.1-x86_64-unknown-linux-musl.tar.gz | tar xz

FROM codegolf/lang-base

COPY --from=0 /bin                     /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/

RUN mkdir -p /gleam/src
COPY gleam.toml /gleam
COPY gleam /usr/bin/

ENTRYPOINT ["gleam"]

CMD ["--version"]
