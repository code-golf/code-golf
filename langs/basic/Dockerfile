FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc libncurses-dev libtinfo5 xz-utils

ENV VER=1.10.1

RUN curl -#L https://downloads.sourceforge.net/fbc/FreeBASIC-$VER-linux-x86_64.tar.xz \
  | tar xJ --strip-components 1

RUN ./install.sh -i /usr

COPY basic.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/basic -s basic.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                        /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                  \
              /lib/x86_64-linux-gnu/libdl.so.2                 \
              /lib/x86_64-linux-gnu/libm.so.6                  \
              /lib/x86_64-linux-gnu/libmvec.so.1               \
              /lib/x86_64-linux-gnu/libncurses.so.6            \
              /lib/x86_64-linux-gnu/libpthread.so.0            \
              /lib/x86_64-linux-gnu/libtinfo.so.5              \
              /lib/x86_64-linux-gnu/libtinfo.so.6              \
              /lib/x86_64-linux-gnu/libz.so.1                  /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                      /lib64/
COPY --from=0 /usr/bin/as                                      \
              /usr/bin/basic                                   \
              /usr/bin/fbc                                     \
              /usr/bin/gcc                                     \
              /usr/bin/ld                                      /usr/bin/
COPY --from=0 /usr/include/freebasic/*.bi                      /usr/include/freebasic/
COPY --from=0 /usr/lib/freebasic/linux-x86_64/fbextra.x        \
              /usr/lib/freebasic/linux-x86_64/fbrt0.o          \
              /usr/lib/freebasic/linux-x86_64/libfb.a          /usr/lib/freebasic/linux-x86_64/
COPY --from=0 /usr/lib/gcc/x86_64-linux-gnu/10/cc1             \
              /usr/lib/gcc/x86_64-linux-gnu/10/crtbegin.o      \
              /usr/lib/gcc/x86_64-linux-gnu/10/crtend.o        \
              /usr/lib/gcc/x86_64-linux-gnu/10/libgcc.a        /usr/lib/gcc/x86_64-linux-gnu/10/
COPY --from=0 /usr/lib/x86_64-linux-gnu/crt1.o                 \
              /usr/lib/x86_64-linux-gnu/crti.o                 \
              /usr/lib/x86_64-linux-gnu/crtn.o                 \
              /usr/lib/x86_64-linux-gnu/libbfd-2.35.2-*.so     \
              /usr/lib/x86_64-linux-gnu/libc.so                \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a       \
              /usr/lib/x86_64-linux-gnu/libctf.so.0            \
              /usr/lib/x86_64-linux-gnu/libdl.so               \
              /usr/lib/x86_64-linux-gnu/libgmp.so.10           \
              /usr/lib/x86_64-linux-gnu/libisl.so.23           \
              /usr/lib/x86_64-linux-gnu/libm.so                \
              /usr/lib/x86_64-linux-gnu/libm-2.31.a            \
              /usr/lib/x86_64-linux-gnu/libmpc.so.3            \
              /usr/lib/x86_64-linux-gnu/libmpfr.so.6           \
              /usr/lib/x86_64-linux-gnu/libmvec.so             \
              /usr/lib/x86_64-linux-gnu/libncurses.so          \
              /usr/lib/x86_64-linux-gnu/libopcodes-2.35.2-*.so \
              /usr/lib/x86_64-linux-gnu/libpthread.so          \
              /usr/lib/x86_64-linux-gnu/libtinfo.so            \
              /usr/lib/x86_64-linux-gnu/libzstd.so.1           /usr/lib/x86_64-linux-gnu/

ENTRYPOINT ["basic"]

CMD ["--version"]
