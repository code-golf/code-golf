FROM codegolf/lang-elixir

FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc patch

ENV VER=c31668b

COPY --from=0 /usr /usr

WORKDIR /usr/local

RUN curl -#L https://github.com/Adriandmen/05AB1E/tarball/$VER \
  | tar xz --strip-components 1

COPY interp.patch /

RUN patch -p1 < /interp.patch

RUN mix local.hex --force          \
 && mix local.rebar                \
 && mix deps.get                   \
 && mix deps.update ssl_verify_fun \
 && MIX_ENV='prod' mix escript.build

COPY 05ab1e.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/05ab1e -s /05ab1e.c

FROM codegolf/lang-base

COPY --from=0 /                 /
COPY --from=1 /usr/bin/05ab1e   /usr/bin/
COPY --from=1 /usr/local/osabie /usr/local/bin/

ENTRYPOINT ["05ab1e"]

CMD ["--version"]
