FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl dash ncurses-dev perl

ENV VER=27.2

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$VER/otp_src_$VER.tar.gz | tar xz

WORKDIR /otp_src_$VER

RUN ./configure && make -j`nproc` install

COPY erlang.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/erlang -s -static /erlang.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                                       /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                            /lib/
COPY --from=0 /usr/bin/basename /usr/bin/dirname /usr/bin/erlang  /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1                              \
              /usr/lib/libncursesw.so.6                           \
              /usr/lib/libstdc++.so.6                             /usr/lib/
COPY --from=0 /usr/local/bin/erl /usr/local/bin/escript           /usr/local/bin/
COPY --from=0 /usr/local/lib/erlang/bin/no_dot_erlang.boot        \
              /usr/local/lib/erlang/bin/start.boot                /usr/local/lib/erlang/bin/
COPY --from=0 /usr/local/lib/erlang/erts-15.2/bin/beam.smp        \
              /usr/local/lib/erlang/erts-15.2/bin/erl_child_setup \
              /usr/local/lib/erlang/erts-15.2/bin/erlexec         \
              /usr/local/lib/erlang/erts-15.2/bin/inet_gethost    /usr/local/lib/erlang/erts-15.2/bin/
COPY --from=0 /usr/local/lib/erlang/lib/compiler-8.5.4/ebin       /usr/local/lib/erlang/lib/compiler-8.5.4/ebin
COPY --from=0 /usr/local/lib/erlang/lib/kernel-10.2/ebin          /usr/local/lib/erlang/lib/kernel-10.2/ebin
COPY --from=0 /usr/local/lib/erlang/lib/stdlib-6.2/ebin           /usr/local/lib/erlang/lib/stdlib-6.2/ebin

ENTRYPOINT ["erlang"]

CMD ["--version"]
