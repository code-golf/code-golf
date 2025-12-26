FROM alpine:3.23 AS builder

RUN apk add --no-cache curl dash gcc musl-dev patch

ENV VER=5bbfddb

WORKDIR /usr/local/bin

RUN curl -#L https://github.com/darrenks/nibbles/tarball/$VER                                     \
  | tar xz --strip-components 1                                                                   \
 && curl -#L https://hackage.haskell.org/package/dlist-1.0/dlist-1.0.tar.gz                       \
  | tar xz --strip-components 1                                                                   \
 && curl -#L https://hackage.haskell.org/package/memoize-1.1.2/memoize-1.1.2.tar.gz               \
  | tar xz --strip-components 2                                                                   \
 && curl -#L https://hackage.haskell.org/package/murmur-hash-0.1.0.11/murmur-hash-0.1.0.11.tar.gz \
  | tar xz --strip-components 1                                                                   \
 && curl -#L https://hackage.haskell.org/package/split-0.2.5/split-0.2.5.tar.gz                   \
  | tar xz --strip-components 2

COPY nibbles.patch .

RUN patch -p0 < nibbles.patch

COPY nibbles.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nibbles -s -static nibbles.c

FROM codegolf/lang-haskell

ENV PATH=/usr/bin:$PATH

COPY --from=0 /usr/bin/dash                              /bin/sh
COPY --from=0 /usr/bin/as /usr/bin/ld /usr/bin/nibbles   /usr/bin/
COPY --from=0 /usr/include                               /usr/include
COPY --from=0 /usr/lib/crt1.o                            \
              /usr/lib/crti.o                            \
              /usr/lib/crtn.o                            \
              /usr/lib/libbfd-2.45.1.so                  \
              /usr/lib/libc.so                           \
              /usr/lib/libctf.so.0                       \
              /usr/lib/libdl.a                           \
              /usr/lib/libgcc_s.so                       \
              /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libisl.so.23                      \
              /usr/lib/libjansson.so.4                   \
              /usr/lib/libm.a                            \
              /usr/lib/libmpc.so.3                       \
              /usr/lib/libmpfr.so.6                      \
              /usr/lib/librt.a                           \
              /usr/lib/libsframe.so.2                    \
              /usr/lib/libssp_nonshared.a                \
              /usr/lib/libz.so.1                         \
              /usr/lib/libzstd.so.1                      /usr/lib/
COPY --from=0 /usr/lib/gcc/*/15.2.0/crtbegin.o           \
              /usr/lib/gcc/*/15.2.0/crtbeginS.o          \
              /usr/lib/gcc/*/15.2.0/crtend.o             \
              /usr/lib/gcc/*/15.2.0/crtendS.o            \
              /usr/lib/gcc/*/15.2.0/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/15.2.0/
COPY --from=0 /usr/libexec/gcc/*/15.2.0/cc1              \
              /usr/libexec/gcc/*/15.2.0/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/15.2.0/
COPY --from=0 /usr/local/bin                             /usr/local/bin

WORKDIR /usr/local/bin

RUN ghc -O -dynamic nibbles.hs

ENTRYPOINT ["nibbles"]
