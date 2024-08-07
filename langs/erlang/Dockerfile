FROM alpine:3.20 AS builder

ARG OTP_VERSION=27.0.1

RUN apk add --no-cache build-base curl ncurses-dev openjdk21=21.0.4_p7-r0 openssl openssl-dev perl=5.38.2-r0

RUN curl -L https://github.com/erlang/otp/releases/download/OTP-$OTP_VERSION/otp_src_$OTP_VERSION.tar.gz | tar xz

RUN mv otp_src_* otp           \
 && cd otp                     \
 && export ERL_TOP=`pwd`       \
 && ./configure                \
    --disable-silent-rules     \
    --prefix=/usr              \
 && export MAKEFLAGS=-j`nproc` \
 && make                       \
 && make install

COPY erlang /usr/bin/
RUN chmod +x /usr/bin/erlang

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                                /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                        /lib/
COPY --from=0 /usr/bin/basename                               \
              /usr/bin/dirname                                \
              /usr/bin/erl                                    \
              /usr/bin/erlang                                 \
              /usr/bin/escript                                /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1                          \
              /usr/lib/libncursesw.so.6                       \
              /usr/lib/libstdc++.so.6                         /usr/lib/
COPY --from=0 /usr/lib/erlang/bin/no_dot_erlang.boot          \
              /usr/lib/erlang/bin/start.boot                  /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-15.0.1/bin/beam.smp        \
              /usr/lib/erlang/erts-15.0.1/bin/erl_child_setup \
              /usr/lib/erlang/erts-15.0.1/bin/erlexec         \
              /usr/lib/erlang/erts-15.0.1/bin/inet_gethost    /usr/lib/erlang/erts-15.0.1/bin/
COPY --from=0 /usr/lib/erlang/lib/compiler-8.5.1/ebin         /usr/lib/erlang/lib/compiler-8.5.1/ebin
COPY --from=0 /usr/lib/erlang/lib/kernel-10.0.1/ebin          /usr/lib/erlang/lib/kernel-10.0.1/ebin
COPY --from=0 /usr/lib/erlang/lib/stdlib-6.0.1/ebin           /usr/lib/erlang/lib/stdlib-6.0.1/ebin

ENTRYPOINT ["erlang"]

CMD ["-version"]
