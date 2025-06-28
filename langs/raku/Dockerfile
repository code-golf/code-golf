FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base git perl

ENV VER=2025.05

RUN git clone --branch $VER https://github.com/rakudo/rakudo.git

WORKDIR /rakudo

RUN CFLAGS='-flto'            \
    ./Configure.pl            \
    --gen-moar                \
    --moar-option=--ar=gcc-ar \
    --prefix=/usr             \
 && make -j`nproc` install    \
 && strip /usr/bin/rakudo     \
          /usr/lib/libmoar.so \
 && rm -r /usr/share/nqp/lib/profiler

WORKDIR /usr/share/perl6/runtime

RUN rm rakudo.moarvm       \
       rakudo-debug.moarvm \
       perl6-debug.moarvm

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/rakudo          /usr/bin/raku
COPY --from=0 /usr/lib/libmoar.so      /usr/lib/
COPY --from=0 /usr/share/nqp           /usr/share/nqp
COPY --from=0 /usr/share/perl6         /usr/share/perl6

ENTRYPOINT ["raku"]

CMD ["--version"]
