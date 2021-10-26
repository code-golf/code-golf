FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base git perl

RUN git clone -b 2021.10 https://github.com/rakudo/rakudo.git \
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

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin/rakudo          /usr/bin/raku
COPY --from=0 /usr/lib/libmoar.so      /usr/lib/
COPY --from=0 /usr/share/nqp           /usr/share/nqp
COPY --from=0 /usr/share/perl6/lib     /usr/share/perl6/lib
COPY --from=0 /usr/share/perl6/runtime /usr/share/perl6/runtime

ENTRYPOINT ["raku"]

CMD ["-v"]
