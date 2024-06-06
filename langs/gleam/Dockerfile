FROM elixir:1.16.3-slim AS builder

RUN apk add --no-cache curl

RUN curl -L https://github.com/gleam-lang/gleam/releases/download/v1.2.1/gleam-v1.2.1-x86_64-unknown-linux-musl.tar.gz | tar xz

FROM codegolf/lang-base

COPY --from=0 /gleam /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/

RUN mkdir -p /tmp/gleam/src
COPY gleam.toml /tmp/gleam/
COPY gleamwrapper /usr/bin/

RUN cd /tmp/gleam && /usr/bin/gleam deps download

ENTRYPOINT ["gleamwrapper"]

CMD ["--version"]
