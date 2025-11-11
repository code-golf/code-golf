FROM codegolf/lang-erlang

FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl make

ENV LANG='C.UTF-8' VER=1.19.2

COPY --from=0 /usr /usr

WORKDIR /usr/local

RUN curl -#L https://github.com/elixir-lang/elixir/archive/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make

FROM codegolf/lang-base

COPY --from=0 /                          /
COPY --from=1 /usr/lib/locale            /usr/lib/locale
COPY --from=1 /usr/local/bin/elixir      /usr/local/bin/
COPY --from=1 /usr/local/lib/elixir/ebin /usr/local/lib/elixir/ebin

ENTRYPOINT ["elixir"]

CMD ["--short-version"]
