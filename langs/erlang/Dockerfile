FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc libssl-dev make openssl perl

# Also rebuild 05AB1E, Elixir, Gleam.
ENV VER=28.1.1

WORKDIR /erlang

RUN curl -#L https://github.com/erlang/otp/releases/download/OTP-$VER/otp_src_$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure             \
    --enable-m64-build      \
    --prefix=/usr           \
    --without-common_test   \
    --without-debugger      \
    --without-dialyzer      \
    --without-diameter      \
    --without-edoc          \
    --without-eldap         \
    --without-et            \
    --without-eunit         \
    --without-ftp           \
    --without-javac         \
    --without-jinterface    \
    --without-megaco        \
    --without-mnesia        \
    --without-observer      \
    --without-odbc          \
    --without-os_mon        \
    --without-reltool       \
    --without-runtime_tools \
    --without-snmp          \
    --without-ssh           \
    --without-termcap       \
    --without-tftp          \
    --without-wx            \
    --without-xmerl         \
 && make -j`nproc` release  \
 && make install            \
 && strip /usr/bin/escript

WORKDIR /usr/lib/erlang/lib

RUN rm -rf erl_interface* erts* sasl*   \
 && mv ../erts* /usr/local              \
 && find /usr/local -type f -not \(     \
    -name 'beam.smp'        -or         \
    -name 'erl_child_setup' -or         \
    -name 'erlexec'         -or         \
    -name 'inet_gethost' \) -delete     \
 && mv /usr/lib/erlang/lib /usr/local   \
 && find /usr/local/lib -type f -not \( \
    -name '*.*a*'           -or         \
    -name 'OTP-PUB-KEY.hrl' -or         \
    -name 'asn1rt_nif.so'   -or         \
    -name 'crypto*.so'      -or         \
    -name 'file.hrl'        -or         \
    -name 'logger.hrl'      -or         \
    -name 'public_key.hrl'  -or         \
    -name 'yecc*' \) -delete

COPY erlang.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/erlang -s /erlang.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                              /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6        \
              /lib/x86_64-linux-gnu/libm.so.6        /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2            /lib64/
COPY --from=0 /usr/bin/basename                      \
              /usr/bin/dirname                       \
              /usr/bin/erl                           \
              /usr/bin/erlang                        \
              /usr/bin/erlc                          \
              /usr/bin/escript                       /usr/bin/
COPY --from=0 /usr/local                             /usr/lib/erlang
COPY --from=0 /usr/lib/erlang/bin/erl                \
              /usr/lib/erlang/bin/no_dot_erlang.boot \
              /usr/lib/erlang/bin/start.boot         /usr/lib/erlang/bin/
COPY --from=0 /usr/lib/locale                        /usr/lib/locale

ENTRYPOINT ["erlang"]

CMD ["--version"]
