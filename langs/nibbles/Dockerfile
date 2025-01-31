ARG GHCVER=8.10.7

FROM alpine:3.21 AS builder

ARG GHCVER

RUN apk add --no-cache build-base curl dash ncurses-dev

ENV VER=f0f01c8

RUN curl -#L https://downloads.haskell.org/~ghc/$GHCVER/ghc-$GHCVER-x86_64-alpine3.10-linux-integer-simple.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure   \
    --prefix=/usr \
 && make install

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

RUN ghc -O -package ghc nibbles.hs \
 && strip /lib/x86_64-linux-ghc-$GHCVER/*.so*

COPY nibbles.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nibbles -s -static /nibbles.c

FROM codegolf/lang-base

ARG GHCVER

COPY --from=0 /usr/bin/dash                                                              /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                                                   /lib/
COPY --from=0 /usr/bin/as /usr/bin/cc /usr/bin/ghc /usr/bin/ld /usr/bin/nibbles          /usr/bin/
COPY --from=0 /usr/include                                                               /usr/include
COPY --from=0 /usr/lib/crt1.o                                                            \
              /usr/lib/crti.o                                                            \
              /usr/lib/crtn.o                                                            \
              /usr/lib/libbfd-2.43.1.so                                                  \
              /usr/lib/libc.so                                                           \
              /usr/lib/libctf.so.0                                                       \
              /usr/lib/libdl.a                                                           \
              /usr/lib/libgcc_s.so                                                       \
              /usr/lib/libgcc_s.so.1                                                     \
              /usr/lib/libgmp.so.10                                                      \
              /usr/lib/libisl.so.23                                                      \
              /usr/lib/libjansson.so.4                                                   \
              /usr/lib/libm.a                                                            \
              /usr/lib/libmpc.so.3                                                       \
              /usr/lib/libmpfr.so.6                                                      \
              /usr/lib/libncursesw.so.6                                                  \
              /usr/lib/librt.a                                                           \
              /usr/lib/libsframe.so.1                                                    \
              /usr/lib/libssp_nonshared.a                                                \
              /usr/lib/libz.so.1                                                         \
              /usr/lib/libzstd.so.1                                                      /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtbegin.o                    \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtbeginS.o                   \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtend.o                      \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtendS.o                     \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/libgcc.a                      /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/lib/ghc-$GHCVER/bin/ghc                                               /usr/lib/ghc-$GHCVER/bin/
COPY --from=0 /usr/lib/ghc-$GHCVER/lib/platformConstants                                 \
              /usr/lib/ghc-$GHCVER/lib/settings                                          /usr/lib/ghc-$GHCVER/lib/
COPY --from=0 /usr/lib/ghc-$GHCVER/lib/package.conf.d                                    /usr/lib/ghc-$GHCVER/lib/package.conf.d
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/libHSarray-0.5.4.0-ghc$GHCVER.so             \
              /lib/x86_64-linux-ghc-$GHCVER/libHSbase-4.14.3.0-ghc$GHCVER.so             \
              /lib/x86_64-linux-ghc-$GHCVER/libHSbinary-0.8.8.0-ghc$GHCVER.so            \
              /lib/x86_64-linux-ghc-$GHCVER/libHSbytestring-0.10.12.0-ghc$GHCVER.so      \
              /lib/x86_64-linux-ghc-$GHCVER/libHScontainers-0.6.5.1-ghc$GHCVER.so        \
              /lib/x86_64-linux-ghc-$GHCVER/libHSdeepseq-1.4.4.0-ghc$GHCVER.so           \
              /lib/x86_64-linux-ghc-$GHCVER/libHSdirectory-1.3.6.0-ghc$GHCVER.so         \
              /lib/x86_64-linux-ghc-$GHCVER/libHSexceptions-0.10.4-ghc$GHCVER.so         \
              /lib/x86_64-linux-ghc-$GHCVER/libHSfilepath-1.4.2.1-ghc$GHCVER.so          \
              /lib/x86_64-linux-ghc-$GHCVER/libHSghc-$GHCVER-ghc$GHCVER.so               \
              /lib/x86_64-linux-ghc-$GHCVER/libHSghc-boot-$GHCVER-ghc$GHCVER.so          \
              /lib/x86_64-linux-ghc-$GHCVER/libHSghc-boot-th-$GHCVER-ghc$GHCVER.so       \
              /lib/x86_64-linux-ghc-$GHCVER/libHSghc-heap-$GHCVER-ghc$GHCVER.so          \
              /lib/x86_64-linux-ghc-$GHCVER/libHSghc-prim-0.6.1-ghc$GHCVER.so            \
              /lib/x86_64-linux-ghc-$GHCVER/libHSghci-$GHCVER-ghc$GHCVER.so              \
              /lib/x86_64-linux-ghc-$GHCVER/libHShaskeline-0.8.2-ghc$GHCVER.so           \
              /lib/x86_64-linux-ghc-$GHCVER/libHShpc-0.6.1.0-ghc$GHCVER.so               \
              /lib/x86_64-linux-ghc-$GHCVER/libHSinteger-simple-0.1.2.0-ghc$GHCVER.so    \
              /lib/x86_64-linux-ghc-$GHCVER/libHSmtl-2.2.2-ghc$GHCVER.so                 \
              /lib/x86_64-linux-ghc-$GHCVER/libHSpretty-1.1.3.6-ghc$GHCVER.so            \
              /lib/x86_64-linux-ghc-$GHCVER/libHSprocess-1.6.13.2-ghc$GHCVER.so          \
              /lib/x86_64-linux-ghc-$GHCVER/libHSrts_thr-ghc$GHCVER.so                   \
              /lib/x86_64-linux-ghc-$GHCVER/libffi.so.7                                  \
              /lib/x86_64-linux-ghc-$GHCVER/libHSrts-1.0.1_thr-ghc$GHCVER.so             \
              /lib/x86_64-linux-ghc-$GHCVER/libHSstm-2.5.0.1-ghc$GHCVER.so               \
              /lib/x86_64-linux-ghc-$GHCVER/libHStemplate-haskell-2.16.0.0-ghc$GHCVER.so \
              /lib/x86_64-linux-ghc-$GHCVER/libHSterminfo-0.4.1.4-ghc$GHCVER.so          \
              /lib/x86_64-linux-ghc-$GHCVER/libHStime-1.9.3-ghc$GHCVER.so                \
              /lib/x86_64-linux-ghc-$GHCVER/libHStransformers-0.5.6.2-ghc$GHCVER.so      \
              /lib/x86_64-linux-ghc-$GHCVER/libHSunix-2.7.2.2-ghc$GHCVER.so              /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/array-0.5.4.0                                /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/array-0.5.4.0
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/base-4.14.3.0                                /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/base-4.14.3.0
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/bytestring-0.10.12.0                         /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/bytestring-0.10.12.0
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/containers-0.6.5.1                           /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/containers-0.6.5.1
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/deepseq-1.4.4.0                              /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/deepseq-1.4.4.0
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/ghc-boot-th-$GHCVER                          /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/ghc-boot-th-$GHCVER
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/ghc-prim-0.6.1                               /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/ghc-prim-0.6.1
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/integer-simple-0.1.2.0                       /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/integer-simple-0.1.2.0
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/pretty-1.1.3.6                               /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/pretty-1.1.3.6
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/rts-1.0.1                                    /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/rts-1.0.1
COPY --from=0 /lib/x86_64-linux-ghc-$GHCVER/template-haskell-2.16.0.0                    /usr/lib/ghc-$GHCVER/lib/x86_64-linux-ghc-$GHCVER/template-haskell-2.16.0.0
COPY --from=0 /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/cc1                       \
              /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/liblto_plugin.so          /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/local/bin/nibbles                                                     /usr/local/bin/
COPY --from=0 /usr/local/bin/Data/Digest/Murmur64.hi                                     \
              /usr/local/bin/Data/Digest/Murmur64.hs                                     \
              /usr/local/bin/Data/Digest/Murmur64.o                                      /usr/local/bin/Data/Digest/
COPY --from=0 /usr/local/bin/Data/Function/Memoize.hi                                    \
              /usr/local/bin/Data/Function/Memoize.hs                                    \
              /usr/local/bin/Data/Function/Memoize.o                                     /usr/local/bin/Data/Function/
COPY --from=0 /usr/local/bin/Data/Function/Memoize/Class.dyn_o                           \
              /usr/local/bin/Data/Function/Memoize/Class.hi                              \
              /usr/local/bin/Data/Function/Memoize/Class.hs                              \
              /usr/local/bin/Data/Function/Memoize/Class.o                               \
              /usr/local/bin/Data/Function/Memoize/TH.dyn_o                              \
              /usr/local/bin/Data/Function/Memoize/TH.hi                                 \
              /usr/local/bin/Data/Function/Memoize/TH.hs                                 \
              /usr/local/bin/Data/Function/Memoize/TH.o                                  /usr/local/bin/Data/Function/Memoize/
COPY --from=0 /usr/local/bin/Data/List/Split.hi                                          \
              /usr/local/bin/Data/List/Split.hs                                          \
              /usr/local/bin/Data/List/Split.o                                           /usr/local/bin/Data/List/
COPY --from=0 /usr/local/bin/Data/List/Split/Internals.hi                                \
              /usr/local/bin/Data/List/Split/Internals.hs                                \
              /usr/local/bin/Data/List/Split/Internals.o                                 /usr/local/bin/Data/List/Split/

ENTRYPOINT ["nibbles"]

CMD ["--version"]
