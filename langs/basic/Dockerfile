FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc libncurses5-dev make xz-utils

ENV VER=1.10.1

WORKDIR /usr/local

RUN curl -#L https://downloads.sourceforge.net/fbc/FreeBASIC-$VER-source-bootstrap.tar.xz \
  | tar xJ --strip-components 1

RUN make -j bootstrap-minimal \
            install-includes  \
 && make -j install-compiler  \
            install-includes  \
            install-rtlib     \
            prefix=/usr

RUN strip /usr/bin/fbc

WORKDIR /usr/include/freebasic
# Remove some unneeded headers.
RUN rm -rf GL X11 gtk win

COPY basic.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/basic -s /basic.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                             /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libbfd-2.44-system.so           \
              /lib/x86_64-linux-gnu/libc.so.6                       \
              /lib/x86_64-linux-gnu/libctf.so.0                     \
              /lib/x86_64-linux-gnu/libgcc_s.so.1                   \
              /lib/x86_64-linux-gnu/libgmp.so.10                    \
              /lib/x86_64-linux-gnu/libisl.so.23                    \
              /lib/x86_64-linux-gnu/libjansson.so.4                 \
              /lib/x86_64-linux-gnu/libm.so.6                       \
              /lib/x86_64-linux-gnu/libmpc.so.3                     \
              /lib/x86_64-linux-gnu/libmpfr.so.6                    \
              /lib/x86_64-linux-gnu/libmvec.so.1                    \
              /lib/x86_64-linux-gnu/libncurses.so.6                 \
              /lib/x86_64-linux-gnu/libsframe.so.1                  \
              /lib/x86_64-linux-gnu/libtinfo.so.6                   \
              /lib/x86_64-linux-gnu/libz.so.1                       \
              /lib/x86_64-linux-gnu/libzstd.so.1                    /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                           /lib64/
COPY --from=0 /usr/bin/as                                           \
              /usr/bin/basic                                        \
              /usr/bin/fbc                                          \
              /usr/bin/gcc                                          \
              /usr/bin/ld                                           /usr/bin/
COPY --from=0 /usr/include/freebasic                                /usr/include/freebasic
COPY --from=0 /usr/lib/freebasic/linux-x86_64                       /usr/lib/freebasic/linux-x86_64
COPY --from=0 /usr/lib/gcc/x86_64-linux-gnu/14/crtbegin.o           \
              /usr/lib/gcc/x86_64-linux-gnu/14/crtbeginS.o          \
              /usr/lib/gcc/x86_64-linux-gnu/14/crtbeginT.o          \
              /usr/lib/gcc/x86_64-linux-gnu/14/crtend.o             \
              /usr/lib/gcc/x86_64-linux-gnu/14/crtendS.o            \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc.a             \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc_eh.a          \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc_s.so          /usr/lib/gcc/x86_64-linux-gnu/14/
COPY --from=0 /usr/lib/x86_64-linux-gnu/Scrt1.o                     \
              /usr/lib/x86_64-linux-gnu/crt1.o                      \
              /usr/lib/x86_64-linux-gnu/crti.o                      \
              /usr/lib/x86_64-linux-gnu/crtn.o                      \
              /usr/lib/x86_64-linux-gnu/libc.a                      \
              /usr/lib/x86_64-linux-gnu/libc.so                     \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a            \
              /usr/lib/x86_64-linux-gnu/libdl.a                     \
              /usr/lib/x86_64-linux-gnu/libm.so                     \
              /usr/lib/x86_64-linux-gnu/libncurses.so               \
              /usr/lib/x86_64-linux-gnu/libpthread.a                \
              /usr/lib/x86_64-linux-gnu/libtinfo.so                 /usr/lib/x86_64-linux-gnu/
COPY --from=0 /usr/libexec/gcc/x86_64-linux-gnu/14/cc1              \
              /usr/libexec/gcc/x86_64-linux-gnu/14/liblto_plugin.so /usr/libexec/gcc/x86_64-linux-gnu/14/

ENTRYPOINT ["basic"]

CMD ["--version"]
