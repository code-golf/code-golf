FROM alpine:3.20 AS builder

ARG OTP_VERSION=OTP-27.0.1

RUN apk add --no-cache curl g++ make ncurses-dev perl

RUN curl -L https://github.com/erlang/otp/archive/refs/tags/$OTP_VERSION.tar.gz | tar xz

ENV ERL_TOP=/otp-$OTP_VERSION

WORKDIR $ERL_TOP

RUN ./configure --disable-jit --prefix=/usr \
 && make --jobs=`nproc` && make --jobs=`nproc` install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/chmod /bin/sh                     /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                        /lib/
COPY          /erlang                                         /usr/bin/
COPY --from=0 /usr/bin/basename /usr/bin/dirname /usr/bin/erl /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6                       /usr/lib/
COPY --from=0 /usr/lib/erlang/bin/start.boot                  /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/erlang/erts-15.0.1/bin/beam.smp        \
              /usr/lib/erlang/erts-15.0.1/bin/erl_child_setup \
              /usr/lib/erlang/erts-15.0.1/bin/erlexec         \
              /usr/lib/erlang/erts-15.0.1/bin/inet_gethost    /usr/lib/erlang/erts-15.0.1/bin/
COPY --from=0 /usr/lib/erlang/lib/kernel-10.0.1/ebin          /usr/lib/erlang/lib/kernel-10.0.1/ebin
COPY --from=0 /usr/lib/erlang/lib/stdlib-6.0.1/ebin           /usr/lib/erlang/lib/stdlib-6.0.1/ebin

RUN chmod +x /usr/bin/erlang

ENTRYPOINT ["erlang"]

CMD ["version"]
