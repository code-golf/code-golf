FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes build-essential curl ncurses-dev

ENV ERLVER=28.0.1 VER=1.18.4

WORKDIR /erlang

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$ERLVER/otp_src_$ERLVER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure   \
    --prefix=/usr \
 && make -j`nproc` install

WORKDIR /elixir

RUN curl -#L https://github.com/elixir-lang/elixir/archive/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN LANG='C.UTF-8' make clean compile

WORKDIR /usr/lib/erlang

RUN mv erts* /usr/local && find /usr/local -type f          \
    -not \( -name 'beam.smp' -or -name 'erl_child_setup'    \
    -or -name 'erlexec' -or -name 'inet_gethost' \) -delete \
 && mv lib/compiler* lib/kernel* lib/stdlib* /usr/local/lib \
 && find /usr/local/lib -type f -not -name '*.*a*' -delete

FROM codegolf/lang-base

ENV ELIXIR_ERL_OPTIONS='+fnu'

COPY --from=0 /bin/dash                            /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/libm.so.6      \
              /lib/x86_64-linux-gnu/libstdc++.so.6 \
              /lib/x86_64-linux-gnu/libtinfo.so.6  /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=0 /usr/bin/basename                    \
              /usr/bin/dirname                     \
              /usr/bin/erl                         /usr/bin/
COPY --from=0 /usr/local                           /usr/lib/erlang
COPY --from=0 /usr/lib/erlang/bin/start.boot       /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/locale                      /usr/lib/locale
COPY --from=0 /elixir/bin/elixir                   /usr/local/bin/
COPY --from=0 /elixir/lib/elixir/ebin              /usr/local/lib/elixir/ebin

ENTRYPOINT ["elixir"]

CMD ["--short-version"]
