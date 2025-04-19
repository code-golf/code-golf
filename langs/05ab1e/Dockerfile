FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes build-essential curl erlang-crypto erlang-dev erlang-ssl ncurses-dev

ENV ERLVER=27.3.3 EXVER=1.18.3 VER=c31668b

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$ERLVER/otp_src_$ERLVER.tar.gz | tar xz

WORKDIR /otp_src_$ERLVER

RUN ./configure   \
    --prefix=/usr \
 && make -j`nproc` install

WORKDIR /usr/local

RUN curl -#L https://github.com/elixir-lang/elixir/archive/v$EXVER.tar.gz \
  | tar xz --strip-components 1

RUN LANG='C.UTF-8' make clean compile

WORKDIR /05ab1e

RUN curl -#L https://github.com/Adriandmen/05AB1E/tarball/$VER \
  | tar xz --strip-components 1

COPY 05ab1e.c interp.patch /

RUN patch -p1 < /interp.patch

RUN mix local.hex --force          \
 && mix local.rebar                \
 && mix deps.get                   \
 && mix deps.update ssl_verify_fun \
 && MIX_ENV='prod' mix escript.build

RUN gcc -Wall -Werror -Wextra -o /usr/bin/05ab1e -s /05ab1e.c

FROM codegolf/lang-base

COPY --from=0 /05ab1e/osabie                                  /
COPY --from=0 /bin/dash                                       /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                 \
              /lib/x86_64-linux-gnu/libgcc_s.so.1             \
              /lib/x86_64-linux-gnu/libm.so.6                 \
              /lib/x86_64-linux-gnu/libstdc++.so.6            \
              /lib/x86_64-linux-gnu/libtinfo.so.6             \
              /lib/x86_64-linux-gnu/libz.so.1                 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                     /lib64/
COPY --from=0 /usr/bin/05ab1e                                 \
              /usr/bin/basename                               \
              /usr/bin/dirname                                \
              /usr/bin/erl                                    \
              /usr/bin/escript                                /usr/bin/
COPY --from=0 /usr/lib/erlang/bin/no_dot_erlang.boot          /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-15.2.6/bin/beam.smp        \
              /usr/lib/erlang/erts-15.2.6/bin/erl_child_setup \
              /usr/lib/erlang/erts-15.2.6/bin/erlexec         \
              /usr/lib/erlang/erts-15.2.6/bin/inet_gethost    /usr/lib/erlang/erts-15.2.6/bin/
COPY --from=0 /usr/lib/erlang/lib/asn1-5.3.4/ebin             /usr/lib/erlang/lib/asn1-5.3.4/ebin
COPY --from=0 /usr/lib/erlang/lib/compiler-8.6.1/ebin         /usr/lib/erlang/lib/compiler-8.6.1/ebin
COPY --from=0 /usr/lib/erlang/lib/crypto-5.1.2/ebin           /usr/lib/erlang/lib/crypto-5.1.2/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-10.2.6/ebin          /usr/lib/erlang/lib/kernel-10.2.6/ebin
COPY --from=0 /usr/lib/erlang/lib/public_key-1.17.1/ebin      /usr/lib/erlang/lib/public_key-1.17.1/ebin
COPY --from=0 /usr/lib/erlang/lib/ssl-10.8.7/ebin             /usr/lib/erlang/lib/ssl-10.8.7/ebin
COPY --from=0 /usr/lib/erlang/lib/stdlib-6.2.2/ebin           /usr/lib/erlang/lib/stdlib-6.2.2/ebin
COPY --from=0 /usr/lib/erlang/lib/syntax_tools-3.2.2/ebin     /usr/lib/erlang/lib/syntax_tools-3.2.2/ebin
COPY --from=0 /usr/lib/erlang/lib/tools-4.1.1/ebin            /usr/lib/erlang/lib/tools-4.1.1/ebin
COPY --from=0 /usr/lib/locale                                 /usr/lib/locale

ENTRYPOINT ["05ab1e"]

CMD ["--version"]
