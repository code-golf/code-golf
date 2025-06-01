FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes build-essential curl ncurses-dev

ENV ERLVER=28.0 VER=1.18.4

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$ERLVER/otp_src_$ERLVER.tar.gz | tar xz

WORKDIR /otp_src_$ERLVER

RUN ./configure   \
    --prefix=/usr \
 && make -j`nproc` install

WORKDIR /usr/local

RUN curl -#L https://github.com/elixir-lang/elixir/archive/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN LANG='C.UTF-8' make clean compile

FROM codegolf/lang-base

ENV ELIXIR_ERL_OPTIONS='+fnu'

COPY --from=0 /bin/dash                                     /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6               \
              /lib/x86_64-linux-gnu/libgcc_s.so.1           \
              /lib/x86_64-linux-gnu/libm.so.6               \
              /lib/x86_64-linux-gnu/libstdc++.so.6          \
              /lib/x86_64-linux-gnu/libtinfo.so.6           /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                   /lib64/
COPY --from=0 /usr/bin/basename                             \
              /usr/bin/dirname                              \
              /usr/bin/erl                                  /usr/bin/
COPY --from=0 /usr/lib/erlang/bin/start.boot                /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-16.0/bin/beam.smp        \
              /usr/lib/erlang/erts-16.0/bin/erl_child_setup \
              /usr/lib/erlang/erts-16.0/bin/erlexec         \
              /usr/lib/erlang/erts-16.0/bin/inet_gethost    /usr/lib/erlang/erts-16.0/bin/
COPY --from=0 /usr/lib/erlang/lib/compiler-9.0/ebin         /usr/lib/erlang/lib/compiler-9.0/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-10.3/ebin          /usr/lib/erlang/lib/kernel-10.3/ebin
COPY --from=0 /usr/lib/erlang/lib/stdlib-7.0/ebin           /usr/lib/erlang/lib/stdlib-7.0/ebin
COPY --from=0 /usr/lib/locale                               /usr/lib/locale
COPY --from=0 /usr/local/bin/elixir                         /usr/local/bin/
COPY --from=0 /usr/local/lib/elixir/ebin                    /usr/local/lib/elixir/ebin

ENTRYPOINT ["elixir"]

CMD ["--short-version"]
