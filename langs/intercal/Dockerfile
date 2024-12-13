FROM alpine:3.21 AS builder

RUN apk add --no-cache bash bison build-base dash flex git

ENV VER=0.33

RUN wget -O- http://www.catb.org/~esr/intercal/intercal-$VER.tar.gz | tar xz

WORKDIR /intercal-$VER

RUN git clone https://github.com/VorpalBlade/cfunge.git

RUN etc/cftoec.sh . cfunge

RUN ./configure   \
    --prefix=/usr \
 && make install  \
 && strip /usr/bin/ick

COPY intercal.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/intercal -s -static /intercal.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                                    /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                         /lib/
COPY --from=0 /usr/bin/as                                      \
              /usr/bin/gcc                                     \
              /usr/bin/ick                                     \
              /usr/bin/intercal                                \
              /usr/bin/ld                                      /usr/bin/
COPY --from=0 /usr/include                                     /usr/include
COPY --from=0 /usr/lib/Scrt1.o                                 \
              /usr/lib/crti.o                                  \
              /usr/lib/crtn.o                                  \
              /usr/lib/libbfd-2.43.1.so                        \
              /usr/lib/libc.so                                 \
              /usr/lib/libctf.so.0                             \
              /usr/lib/libgcc_s.so                             \
              /usr/lib/libgcc_s.so.1                           \
              /usr/lib/libgmp.so.10                            \
              /usr/lib/libick.a                                \
              /usr/lib/libisl.so.23                            \
              /usr/lib/libjansson.so.4                         \
              /usr/lib/libmpc.so.3                             \
              /usr/lib/libmpfr.so.6                            \
              /usr/lib/libsframe.so.1                          \
              /usr/lib/libssp_nonshared.a                      \
              /usr/lib/libz.so.1                               \
              /usr/lib/libzstd.so.1                            /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64*/14.2.0/crtbeginS.o          \
              /usr/lib/gcc/x86_64*/14.2.0/crtendS.o            \
              /usr/lib/gcc/x86_64*/14.2.0/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/libexec/gcc/x86_64*/14.2.0/cc1              \
              /usr/libexec/gcc/x86_64*/14.2.0/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/share/ick-0.33/ick-wrap.c                   /usr/share/ick-0.33/

ENTRYPOINT ["intercal"]

CMD ["--version"]
