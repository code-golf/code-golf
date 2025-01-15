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

COPY --from=0 /05ab1e/osabie                                    /
COPY --from=0 /bin/dash                                         /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                   \
              /lib/x86_64-linux-gnu/libgcc_s.so.1               \
              /lib/x86_64-linux-gnu/libm.so.6                   \
              /lib/x86_64-linux-gnu/libstdc++.so.6              \
              /lib/x86_64-linux-gnu/libtinfo.so.6               \
              /lib/x86_64-linux-gnu/libz.so.1                   /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                       /lib64/
COPY --from=0 /usr/bin/05ab1e                                   \
              /usr/bin/basename                                 \
              /usr/bin/dirname                                  \
              /usr/bin/erl                                      \
              /usr/bin/escript                                  /usr/bin/
COPY --from=0 /usr/lib/erlang/bin/no_dot_erlang.boot            /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-13.2.2.9/bin/beam.smp        \
              /usr/lib/erlang/erts-13.2.2.9/bin/erl_child_setup \
              /usr/lib/erlang/erts-13.2.2.9/bin/erlexec         \
              /usr/lib/erlang/erts-13.2.2.9/bin/inet_gethost    /usr/lib/erlang/erts-13.2.2.9/bin/
COPY --from=0 /usr/lib/erlang/lib/asn1-5.0.21.1/ebin            /usr/lib/erlang/lib/asn1-5.0.21.1/ebin
COPY --from=0 /usr/lib/erlang/lib/compiler-8.2.6.4/ebin         /usr/lib/erlang/lib/compiler-8.2.6.4/ebin
COPY --from=0 /usr/lib/erlang/lib/crypto-5.1.4.3/ebin           /usr/lib/erlang/lib/crypto-5.1.4.3/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-8.5.4.3/ebin           /usr/lib/erlang/lib/kernel-8.5.4.3/ebin
COPY --from=0 /usr/lib/erlang/lib/public_key-1.13.3.2/ebin      /usr/lib/erlang/lib/public_key-1.13.3.2/ebin
COPY --from=0 /usr/lib/erlang/lib/ssl-10.9.1.4/ebin             /usr/lib/erlang/lib/ssl-10.9.1.4/ebin
COPY --from=0 /usr/lib/erlang/lib/stdlib-4.3.1.4/ebin           /usr/lib/erlang/lib/stdlib-4.3.1.4/ebin
COPY --from=0 /usr/lib/erlang/lib/syntax_tools-3.0.1/ebin       /usr/lib/erlang/lib/syntax_tools-3.0.1/ebin
COPY --from=0 /usr/lib/erlang/lib/tools-3.5.3/ebin              /usr/lib/erlang/lib/tools-3.5.3/ebin

ENTRYPOINT ["05ab1e"]

CMD ["--version"]
