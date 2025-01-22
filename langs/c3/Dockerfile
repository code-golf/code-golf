FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl dash gcc libxml2-dev

ENV VER=0.6.6

RUN curl -#L https://github.com/c3lang/c3c/releases/download/v$VER/c3-linux.tar.gz \
  | tar xz --directory /usr/local/bin --strip-components 1

COPY c3.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/c3 -s c3.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                         /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                   \
              /lib/x86_64-linux-gnu/libicudata.so.72            \
              /lib/x86_64-linux-gnu/libicuuc.so.72              \
              /lib/x86_64-linux-gnu/liblzma.so.5                \
              /lib/x86_64-linux-gnu/libm.so.6                   \
              /lib/x86_64-linux-gnu/libmvec.so.1                \
              /lib/x86_64-linux-gnu/libstdc++.so.6              \
              /lib/x86_64-linux-gnu/libtinfo.so.6               \
              /lib/x86_64-linux-gnu/libxml2.so.2                \
              /lib/x86_64-linux-gnu/libz.so.1                   \
              /lib/x86_64-linux-gnu/libzstd.so.1                /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                       /lib64/
COPY --from=0 /usr/bin/c3 /usr/bin/cc /usr/bin/ld               /usr/bin/
COPY --from=0 /usr/lib/gcc/x86_64-linux-gnu/12/crtbeginS.o      \
              /usr/lib/gcc/x86_64-linux-gnu/12/crtendS.o        \
              /usr/lib/gcc/x86_64-linux-gnu/12/libgcc.a         \
              /usr/lib/gcc/x86_64-linux-gnu/12/libgcc_s.so      \
              /usr/lib/gcc/x86_64-linux-gnu/12/liblto_plugin.so /usr/lib/gcc/x86_64-linux-gnu/12/
COPY --from=0 /usr/lib/x86_64-linux-gnu/Scrt1.o                 \
              /usr/lib/x86_64-linux-gnu/crti.o                  \
              /usr/lib/x86_64-linux-gnu/crtn.o                  \
              /usr/lib/x86_64-linux-gnu/libbfd-2.40-system.so   \
              /usr/lib/x86_64-linux-gnu/libc.so                 \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a        \
              /usr/lib/x86_64-linux-gnu/libctf.so.0             \
              /usr/lib/x86_64-linux-gnu/libdl.a                 \
              /usr/lib/x86_64-linux-gnu/libgcc_s.so.1           \
              /usr/lib/x86_64-linux-gnu/libjansson.so.4         \
              /usr/lib/x86_64-linux-gnu/libm.so                 \
              /usr/lib/x86_64-linux-gnu/libpthread.a            \
              /usr/lib/x86_64-linux-gnu/libsframe.so.0          /usr/lib/x86_64-linux-gnu/
COPY --from=0 /usr/local/bin/c3c                                /usr/local/bin/
COPY --from=0 /usr/local/bin/lib/std                            /usr/local/bin/lib/std

ENTRYPOINT ["c3"]

CMD ["--version"]
