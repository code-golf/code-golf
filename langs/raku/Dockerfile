FROM alpine:3.15 as builder

RUN apk add --no-cache build-base git perl

RUN git clone -b 2022.04 https://github.com/rakudo/rakudo.git \
 && cd rakudo                                                 \
 && CFLAGS=-flto ./Configure.pl                               \
    --gen-moar                                                \
    --moar-option=--ar=gcc-ar                                 \
    --prefix=/usr                                             \
 && make -j`nproc` install                                    \
 && strip /usr/bin/rakudo /usr/lib/libmoar.so

RUN rm -r /usr/share/nqp/lib/profiler            \
    /usr/share/perl6/runtime/rakudo-debug.moarvm \
    /usr/share/perl6/runtime/rakudo.moarvm       \
    /usr/share/perl6/runtime/perl6-debug.moarvm

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/rakudo          /usr/bin/raku
COPY --from=0 /usr/lib/libmoar.so      /usr/lib/
COPY --from=0 /usr/share/nqp           /usr/share/nqp
COPY --from=0 /usr/share/perl6         /usr/share/perl6

ENTRYPOINT ["raku"]

CMD ["-v"]
