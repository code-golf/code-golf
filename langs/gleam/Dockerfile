FROM codegolf/lang-erlang

FROM rust:1.90.0-slim-trixie AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl make

ENV VER=1.12.0

WORKDIR /gleam

RUN curl -#L https://github.com/gleam-lang/gleam/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install

WORKDIR /usr/local

RUN gleam new --name main . \
 && gleam add argv          \
 && gleam remove gleeunit

COPY gleam.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/gleam -s /gleam.c

FROM codegolf/lang-base

COPY --from=0 /                                   /
COPY --from=1 /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/
COPY --from=1 /usr/bin/gleam                      /usr/bin/
COPY --from=1 /usr/lib/locale                     /usr/lib/locale
COPY --chown=nobody                               \
     --from=1 /usr/local/gleam.toml               \
              /usr/local/manifest.toml            /usr/local/
COPY --chown=nobody                               \
     --from=1 /usr/local/cargo/bin/gleam          /usr/local/bin/
COPY --chown=nobody                               \
     --from=1 /usr/local/build/packages           /usr/local/build/packages
COPY --chown=nobody                               \
     --from=1 /usr/local/src                      /usr/local/src

ENTRYPOINT ["gleam"]

CMD ["--version"]
