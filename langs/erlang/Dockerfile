FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc make ncurses-dev perl

# Also rebuild Elixir, Gleam.
ENV VER=28.0.1

WORKDIR /erlang

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$VER/otp_src_$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure            \
    --prefix=/usr          \
 && make -j`nproc` install \
 && strip /usr/bin/escript

#WORKDIR /usr/lib/erlang
#
#RUN mv erts* /usr/local && find /usr/local -type f          \
#    -not \( -name 'beam.smp' -or -name 'erl_child_setup'    \
#    -or -name 'erlexec' -or -name 'inet_gethost' \) -delete \
# && mv lib/compiler* lib/kernel* lib/stdlib* /usr/local/lib \
# && find /usr/local/lib -type f -not -name '*.beam' -delete

COPY erlang.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/erlang -s /erlang.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                               /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                         \
              /lib/x86_64-linux-gnu/libgcc_s.so.1                     \
              /lib/x86_64-linux-gnu/libm.so.6                         \
              /lib/x86_64-linux-gnu/libncursesw.so.6                  \
              /lib/x86_64-linux-gnu/libstdc++.so.6                    \
              /lib/x86_64-linux-gnu/libtinfo.so.6                     /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                             /lib64/
COPY --from=0 /usr/bin/basename                                       \
              /usr/bin/dirname                                        \
              /usr/bin/erl                                            \
              /usr/bin/erlang                                         \
              /usr/bin/erlc                                           \
              /usr/bin/escript                                        /usr/bin/
COPY --from=0 /usr/lib/erlang/bin/no_dot_erlang.boot                  \
              /usr/lib/erlang/bin/start.boot                          /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-16.0.1/bin/beam.smp                \
              /usr/lib/erlang/erts-16.0.1/bin/erl_child_setup         \
              /usr/lib/erlang/erts-16.0.1/bin/erlexec                 \
              /usr/lib/erlang/erts-16.0.1/bin/inet_gethost            /usr/lib/erlang/erts-16.0.1/bin/
COPY --from=0 /usr/lib/erlang/lib/compiler-9.0/ebin                   /usr/lib/erlang/lib/compiler-9.0/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-10.3.1/ebin                  /usr/lib/erlang/lib/kernel-10.3.1/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-10.3.1/include/file.hrl      /usr/lib/erlang/lib/kernel-10.3.1/include/
COPY --from=0 /usr/lib/erlang/lib/parsetools-2.7/ebin/yecc.beam       \
              /usr/lib/erlang/lib/parsetools-2.7/ebin/yeccparser.beam \
              /usr/lib/erlang/lib/parsetools-2.7/ebin/yeccscan.beam   /usr/lib/erlang/lib/parsetools-2.7/ebin/
COPY --from=0 /usr/lib/erlang/lib/parsetools-2.7/include/yeccpre.hrl  /usr/lib/erlang/lib/parsetools-2.7/include/
COPY --from=0 /usr/lib/erlang/lib/stdlib-7.0.1/ebin                   /usr/lib/erlang/lib/stdlib-7.0.1/ebin
COPY --from=0 /usr/lib/erlang/lib/tools-4.1.2/ebin/make.beam          /usr/lib/erlang/lib/tools-4.1.2/ebin/

ENTRYPOINT ["erlang"]

CMD ["--version"]
