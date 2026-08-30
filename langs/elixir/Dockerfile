FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes make

ENV LANG='C.UTF-8' VER=1.20.4

COPY --from=codegolf/lang-erlang /usr /usr

WORKDIR /elixir

ADD https://github.com/elixir-lang/elixir.git#v$VER .

RUN make

FROM codegolf/lang-erlang

COPY --from=0 /elixir/bin/elixir      /usr/bin/
COPY --from=0 /elixir/lib/elixir/ebin /usr/lib/elixir/ebin

ENTRYPOINT ["elixir"]

CMD ["--short-version"]
