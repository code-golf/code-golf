FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl elixir erlang gcc patch

ENV VER=c31668b

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

COPY --from=0 /05ab1e/osabie                               /
COPY --from=0 /bin/dash                                    /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6              \
              /lib/x86_64-linux-gnu/libgcc_s.so.1          \
              /lib/x86_64-linux-gnu/libm.so.6              \
              /lib/x86_64-linux-gnu/libstdc++.so.6         \
              /lib/x86_64-linux-gnu/libtinfo.so.6          \
              /lib/x86_64-linux-gnu/libz.so.1              /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                  /lib64/
COPY --from=0 /usr/bin/05ab1e                              \
              /usr/bin/basename                            \
              /usr/bin/dirname                             \
              /usr/bin/erl                                 \
              /usr/bin/escript                             /usr/bin/
COPY --parents                                             \
     --from=0 /usr/lib/./erlang/bin/no_dot_erlang.boot     \
              /usr/lib/./erlang/erts-*/bin/beam.smp        \
              /usr/lib/./erlang/erts-*/bin/erl_child_setup \
              /usr/lib/./erlang/erts-*/bin/erlexec         \
              /usr/lib/./erlang/erts-*/bin/inet_gethost    \
              /usr/lib/./erlang/lib/asn1-*/ebin            \
              /usr/lib/./erlang/lib/compiler-*/ebin        \
              /usr/lib/./erlang/lib/crypto-*/ebin          \
              /usr/lib/./erlang/lib/kernel-*/ebin          \
              /usr/lib/./erlang/lib/public_key-*/ebin      \
              /usr/lib/./erlang/lib/ssl-*/ebin             \
              /usr/lib/./erlang/lib/stdlib-*/ebin          \
              /usr/lib/./erlang/lib/syntax_tools-*/ebin    \
              /usr/lib/./erlang/lib/tools-*/ebin           /usr/lib/

ENTRYPOINT ["05ab1e"]

CMD ["--version"]
