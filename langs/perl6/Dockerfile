FROM alpine:3.10 as builder

RUN apk add --no-cache gcc git linux-headers make musl-dev perl

RUN git clone -b 2019.07.1 git://github.com/rakudo/rakudo \
 && cd rakudo                                             \
 && CFLAGS=-flto ./Configure.pl                           \
    --gen-moar                                            \
    --moar-option=--ar=gcc-ar                             \
    --prefix=/usr                                         \
 && make -j`nproc` install                                \
 && strip /usr/bin/perl6

RUN rm -r /usr/share/nqp/lib/profiler \
   /usr/share/perl6/runtime/perl6-debug.moarvm

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/perl6           /usr/bin/
COPY --from=0 /usr/lib/libmoar.so      /usr/lib/
COPY --from=0 /usr/share/nqp           /usr/share/nqp
COPY --from=0 /usr/share/perl6/lib     /usr/share/perl6/lib
COPY --from=0 /usr/share/perl6/runtime /usr/share/perl6/runtime

ENTRYPOINT ["/usr/bin/perl6", "-v"]
