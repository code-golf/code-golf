FROM alpine:edge

RUN apk add --no-cache gcc git linux-headers make musl-dev perl

RUN git config --global user.email 'foo@bar.com' \
 && git config --global user.name  'foo bar'

# MoarVM
RUN git clone git://github.com/MoarVM/MoarVM \
 && cd MoarVM                                \
 && git checkout f59ac65                     \
 && git revert 4ce52ef                       \
 && ./Configure.pl --prefix=/usr             \
 && make -j`nproc` install

# NQP
RUN git clone git://github.com/perl6/nqp \
 && cd nqp                               \
 && git checkout a1efba5                 \
 && ./Configure.pl --prefix=/usr         \
 && make -j`nproc` install

# Rakudo
RUN git clone git://github.com/rakudo/rakudo \
 && cd rakudo                                \
 && git checkout 8019c15                     \
 && ./Configure.pl --prefix=/usr             \
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
