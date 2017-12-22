FROM alpine:edge

RUN apk add --no-cache curl gcc git linux-headers make musl-dev perl

RUN curl -L https://github.com/rakudo/rakudo/archive/2017.12.tar.gz \
   | tar xzf -

RUN cd rakudo-2017.12                                         \
 && perl Configure.pl --backend=moar --gen-moar --prefix=/usr \
 && make -j`nproc` install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/moar            /usr/bin/
COPY --from=0 /usr/lib/libmoar.so      /usr/lib/
COPY --from=0 /usr/share/nqp           /usr/share/nqp
COPY --from=0 /usr/share/perl6         /usr/share/perl6

ENTRYPOINT [                                 \
    "/usr/bin/moar",                         \
    "--libpath=/usr/share/nqp/lib",          \
    "--libpath=/usr/share/perl6/runtime",    \
    "/usr/share/perl6/runtime/perl6.moarvm", \
    "-v"                                     \
]
