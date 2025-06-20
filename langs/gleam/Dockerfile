FROM codegolf/lang-erlang AS builder

FROM rust:1.87.0-slim-bookworm

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl make

ENV VER=1.11.1

WORKDIR /gleam

RUN curl -#L https://github.com/gleam-lang/gleam/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install

WORKDIR /usr/local

RUN gleam new --name main --skip-git .

RUN gleam add argv \
 && gleam remove gleeunit

RUN chmod -R 777 .

COPY gleam.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/gleam -s /gleam.c

FROM codegolf/lang-base

COPY --from=0 /                          /
COPY --from=1 /usr/bin/gleam             /usr/bin/
COPY --from=1 /usr/lib/locale            /usr/lib/locale
COPY --from=1 /usr/local/gleam.toml      \
              /usr/local/manifest.toml   /usr/local/
COPY --from=1 /usr/local/cargo/bin/gleam /usr/local/bin/
COPY --from=1 /usr/local/build/packages  /usr/local/build/packages
COPY --from=1 /usr/local/src             /usr/local/src

ENTRYPOINT ["gleam"]

CMD ["--version"]
