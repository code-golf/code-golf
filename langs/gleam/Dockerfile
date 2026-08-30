FROM rust:1.98.0-slim-trixie AS builder

ENV VER=1.18.1

WORKDIR /gleam

ADD https://github.com/gleam-lang/gleam.git#v$VER .

RUN cargo install --path gleam-bin --force --locked

WORKDIR /usr/local

RUN gleam new --name main . \
 && gleam add argv          \
 && gleam remove gleeunit

RUN strip cargo/bin/gleam

COPY gleam.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/gleam -s gleam.c

FROM codegolf/lang-erlang

COPY --from=0 /lib/x86_64-linux-gnu/libgcc_s.so.1 /lib/
COPY --from=0 /usr/bin/gleam                      /usr/bin/
COPY --from=0 --chmod=755 --parents               \
              /usr/local/./gleam.toml             \
              /usr/local/./manifest.toml          \
              /usr/local/./build/packages         \
              /usr/local/./src                    \
              /usr/local/cargo/./bin/gleam        /usr/local/

WORKDIR /usr/local

RUN gleam build

ENTRYPOINT ["gleam"]

CMD ["--version"]
