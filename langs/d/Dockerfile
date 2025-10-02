FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl dpkg gcc libc6-dev

ENV VER=2.111.0

RUN curl -#o dmd.deb https://downloads.dlang.org/releases/2.x/$VER/dmd_$VER-0_amd64.deb \
 && dpkg -i dmd.deb

COPY d.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/D -s -static d.c

FROM codegolf/lang-base

COPY --from=0 /etc/dmd.conf                                     /etc/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                   \
              /lib/x86_64-linux-gnu/libdl.so.2                  \
              /lib/x86_64-linux-gnu/libgcc_s.so.1               \
              /lib/x86_64-linux-gnu/libm.so.6                   \
              /lib/x86_64-linux-gnu/libmvec.so.1                \
              /lib/x86_64-linux-gnu/libpthread.so.0             \
              /lib/x86_64-linux-gnu/libsframe.so.1              /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                       /lib64/
COPY --from=0 /usr/bin/D                                        \
              /usr/bin/cc                                       \
              /usr/bin/dmd                                      \
              /usr/bin/ld                                       /usr/bin/
COPY --from=0 /usr/include/dmd                                  /usr/include/dmd
COPY --from=0 /usr/lib/gcc/x86_64-linux-gnu/14/crtbeginS.o      \
              /usr/lib/gcc/x86_64-linux-gnu/14/crtendS.o        \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc.a         \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc_s.so      \
              /usr/lib/gcc/x86_64-linux-gnu/14/liblto_plugin.so /usr/lib/gcc/x86_64-linux-gnu/14/
COPY --from=0 /usr/lib/x86_64-linux-gnu/Scrt1.o                 \
              /usr/lib/x86_64-linux-gnu/crti.o                  \
              /usr/lib/x86_64-linux-gnu/crtn.o                  \
              /usr/lib/x86_64-linux-gnu/libbfd-2.44-system.so   \
              /usr/lib/x86_64-linux-gnu/libc.so                 \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a        \
              /usr/lib/x86_64-linux-gnu/libctf.so.0             \
              /usr/lib/x86_64-linux-gnu/libdl.a                 \
              /usr/lib/x86_64-linux-gnu/libgcc_s.so.1           \
              /usr/lib/x86_64-linux-gnu/libjansson.so.4         \
              /usr/lib/x86_64-linux-gnu/libm.so                 \
              /usr/lib/x86_64-linux-gnu/libphobos2.a            \
              /usr/lib/x86_64-linux-gnu/libpthread.a            \
              /usr/lib/x86_64-linux-gnu/librt.a                 \
              /usr/lib/x86_64-linux-gnu/libz.so.1               \
              /usr/lib/x86_64-linux-gnu/libzstd.so.1            /usr/lib/x86_64-linux-gnu/

ENTRYPOINT ["D"]

CMD ["--version"]
