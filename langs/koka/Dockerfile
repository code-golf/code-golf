FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl dash gcc

ENV VER=3.1.2

RUN curl -#L https://github.com/koka-lang/koka/releases/download/v$VER/koka-v$VER-linux-x64.tar.gz \
  | tar xz --directory /usr/local

COPY kokawrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/kokawrapper -s -static kokawrapper.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                                       /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                     \
              /lib/x86_64-linux-gnu/libdl.so.2                    \
              /lib/x86_64-linux-gnu/libm.so.6                     \
              /lib/x86_64-linux-gnu/libmvec.so.1                  \
              /lib/x86_64-linux-gnu/libpthread.so.0               \
              /lib/x86_64-linux-gnu/librt.so.1                    \
              /lib/x86_64-linux-gnu/libutil.so.1                  \
              /lib/x86_64-linux-gnu/libz.so.1                     /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                         /lib64/
COPY --from=0 /usr/bin/as                                         \
              /usr/bin/gcc                                        \
              /usr/bin/kokawrapper                                \
              /usr/bin/ld                                         /usr/bin/
COPY --from=0 /usr/include                                        /usr/include
COPY --from=0 /usr/lib/gcc                                        /usr/lib/gcc
COPY --from=0 /usr/lib/x86_64-linux-gnu/Scrt1.o                   \
              /usr/lib/x86_64-linux-gnu/crti.o                    \
              /usr/lib/x86_64-linux-gnu/crtn.o                    \
              /usr/lib/x86_64-linux-gnu/libbfd-2.40-system.so     \
              /usr/lib/x86_64-linux-gnu/libbrotlicommon.so.1      \
              /usr/lib/x86_64-linux-gnu/libbrotlidec.so.1         \
              /usr/lib/x86_64-linux-gnu/libbrotlienc.so.1         \
              /usr/lib/x86_64-linux-gnu/libc.so                   \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a          \
              /usr/lib/x86_64-linux-gnu/libcrypto.so.3            \
              /usr/lib/x86_64-linux-gnu/libctf.so.0               \
              /usr/lib/x86_64-linux-gnu/libgcc_s.so.1             \
              /usr/lib/x86_64-linux-gnu/libgmp.so.10              \
              /usr/lib/x86_64-linux-gnu/libisl.so.23              \
              /usr/lib/x86_64-linux-gnu/libjansson.so.4           \
              /usr/lib/x86_64-linux-gnu/libm.so                   \
              /usr/lib/x86_64-linux-gnu/libm-2.36.a               \
              /usr/lib/x86_64-linux-gnu/libmpc.so.3               \
              /usr/lib/x86_64-linux-gnu/libmpfr.so.6              \
              /usr/lib/x86_64-linux-gnu/libmvec.a                 \
              /usr/lib/x86_64-linux-gnu/libnghttp2.so.14          \
              /usr/lib/x86_64-linux-gnu/libopcodes-2.40-system.so \
              /usr/lib/x86_64-linux-gnu/libpthread.a              \
              /usr/lib/x86_64-linux-gnu/libsframe.so.0            \
              /usr/lib/x86_64-linux-gnu/libssl.so.3               \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6            \
              /usr/lib/x86_64-linux-gnu/libzstd.so.1              /usr/lib/x86_64-linux-gnu/
COPY --from=0 /usr/local/bin/koka                                 /usr/local/bin/
COPY --from=0 /usr/local/lib/koka                                 /usr/local/lib/koka
COPY --from=0 /usr/local/share/koka                               /usr/local/share/koka

ENTRYPOINT ["kokawrapper"]

CMD ["--version"]
