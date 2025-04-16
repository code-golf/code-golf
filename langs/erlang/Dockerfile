FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes build-essential curl ncurses-dev

# Also bump Elixir, Gleam.
ENV VER=27.3.3

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$VER/otp_src_$VER.tar.gz | tar xz

WORKDIR /otp_src_$VER

RUN ./configure            \
    --prefix=/usr          \
 && make -j`nproc` install \
 && strip /usr/bin/escript

COPY erlang.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/erlang -s /erlang.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                       /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                 \
              /lib/x86_64-linux-gnu/libgcc_s.so.1             \
              /lib/x86_64-linux-gnu/libm.so.6                 \
              /lib/x86_64-linux-gnu/libncursesw.so.6          \
              /lib/x86_64-linux-gnu/libstdc++.so.6            \
              /lib/x86_64-linux-gnu/libtinfo.so.6             /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                     /lib64/
COPY --from=0 /usr/bin/basename                               \
              /usr/bin/dirname                                \
              /usr/bin/erl                                    \
              /usr/bin/erlang                                 \
              /usr/bin/erlc                                   \
              /usr/bin/escript                                /usr/bin/
COPY --from=0 /usr/lib/erlang/bin/no_dot_erlang.boot          \
              /usr/lib/erlang/bin/start.boot                  /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-15.2.6/bin/beam.smp        \
              /usr/lib/erlang/erts-15.2.6/bin/erl_child_setup \
              /usr/lib/erlang/erts-15.2.6/bin/erlexec         \
              /usr/lib/erlang/erts-15.2.6/bin/inet_gethost    /usr/lib/erlang/erts-15.2.6/bin/
COPY --from=0 /usr/lib/erlang/lib/compiler-8.6.1/ebin         /usr/lib/erlang/lib/compiler-8.6.1/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-10.2.6/ebin          /usr/lib/erlang/lib/kernel-10.2.6/ebin
COPY --from=0 /usr/lib/erlang/lib/stdlib-6.2.2/ebin           /usr/lib/erlang/lib/stdlib-6.2.2/ebin

ENTRYPOINT ["erlang"]

CMD ["--version"]
