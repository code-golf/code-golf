FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl dash ncurses-dev

ENV VER=f0f01c8

RUN curl -#L https://downloads.haskell.org/~ghc/8.10.7/ghc-8.10.7-x86_64-alpine3.10-linux-integer-simple.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure  \
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
 && strip /lib/x86_64-linux-ghc-8.10.7/*.so*

COPY nibbleswrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nibbleswrapper -s -static /nibbleswrapper.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                                                            /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                                                 /lib/
COPY --from=0 /usr/bin/as /usr/bin/cc /usr/bin/ghc /usr/bin/ld /usr/bin/nibbleswrapper /usr/bin/
COPY --from=0 /usr/include                                                             /usr/include
COPY --from=0 /usr/lib/crt1.o                                                          \
              /usr/lib/crti.o                                                          \
              /usr/lib/crtn.o                                                          \
              /usr/lib/libbfd-2.43.1.so                                                \
              /usr/lib/libc.so                                                         \
              /usr/lib/libctf.so.0                                                     \
              /usr/lib/libdl.a                                                         \
              /usr/lib/libgcc_s.so                                                     \
              /usr/lib/libgcc_s.so.1                                                   \
              /usr/lib/libgmp.so.10                                                    \
              /usr/lib/libisl.so.23                                                    \
              /usr/lib/libjansson.so.4                                                 \
              /usr/lib/libm.a                                                          \
              /usr/lib/libmpc.so.3                                                     \
              /usr/lib/libmpfr.so.6                                                    \
              /usr/lib/libncursesw.so.6                                                \
              /usr/lib/librt.a                                                         \
              /usr/lib/libsframe.so.1                                                  \
              /usr/lib/libssp_nonshared.a                                              \
              /usr/lib/libz.so.1                                                       \
              /usr/lib/libzstd.so.1                                                    /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtbegin.o                  \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtbeginS.o                 \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtend.o                    \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/crtendS.o                   \
              /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/libgcc.a                    /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/lib/ghc-8.10.7/bin/ghc                                              /usr/lib/ghc-8.10.7/bin/
COPY --from=0 /usr/lib/ghc-8.10.7/lib/platformConstants                                \
              /usr/lib/ghc-8.10.7/lib/settings                                         /usr/lib/ghc-8.10.7/lib/
COPY --from=0 /usr/lib/ghc-8.10.7/lib/package.conf.d                                   /usr/lib/ghc-8.10.7/lib/package.conf.d
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/libHSarray-0.5.4.0-ghc8.10.7.so             \
              /lib/x86_64-linux-ghc-8.10.7/libHSbase-4.14.3.0-ghc8.10.7.so             \
              /lib/x86_64-linux-ghc-8.10.7/libHSbinary-0.8.8.0-ghc8.10.7.so            \
              /lib/x86_64-linux-ghc-8.10.7/libHSbytestring-0.10.12.0-ghc8.10.7.so      \
              /lib/x86_64-linux-ghc-8.10.7/libHScontainers-0.6.5.1-ghc8.10.7.so        \
              /lib/x86_64-linux-ghc-8.10.7/libHSdeepseq-1.4.4.0-ghc8.10.7.so           \
              /lib/x86_64-linux-ghc-8.10.7/libHSdirectory-1.3.6.0-ghc8.10.7.so         \
              /lib/x86_64-linux-ghc-8.10.7/libHSexceptions-0.10.4-ghc8.10.7.so         \
              /lib/x86_64-linux-ghc-8.10.7/libHSfilepath-1.4.2.1-ghc8.10.7.so          \
              /lib/x86_64-linux-ghc-8.10.7/libHSghc-8.10.7-ghc8.10.7.so                \
              /lib/x86_64-linux-ghc-8.10.7/libHSghc-boot-8.10.7-ghc8.10.7.so           \
              /lib/x86_64-linux-ghc-8.10.7/libHSghc-boot-th-8.10.7-ghc8.10.7.so        \
              /lib/x86_64-linux-ghc-8.10.7/libHSghc-heap-8.10.7-ghc8.10.7.so           \
              /lib/x86_64-linux-ghc-8.10.7/libHSghc-prim-0.6.1-ghc8.10.7.so            \
              /lib/x86_64-linux-ghc-8.10.7/libHSghci-8.10.7-ghc8.10.7.so               \
              /lib/x86_64-linux-ghc-8.10.7/libHShaskeline-0.8.2-ghc8.10.7.so           \
              /lib/x86_64-linux-ghc-8.10.7/libHShpc-0.6.1.0-ghc8.10.7.so               \
              /lib/x86_64-linux-ghc-8.10.7/libHSinteger-simple-0.1.2.0-ghc8.10.7.so    \
              /lib/x86_64-linux-ghc-8.10.7/libHSmtl-2.2.2-ghc8.10.7.so                 \
              /lib/x86_64-linux-ghc-8.10.7/libHSpretty-1.1.3.6-ghc8.10.7.so            \
              /lib/x86_64-linux-ghc-8.10.7/libHSprocess-1.6.13.2-ghc8.10.7.so          \
              /lib/x86_64-linux-ghc-8.10.7/libHSrts_thr-ghc8.10.7.so                   \
              /lib/x86_64-linux-ghc-8.10.7/libffi.so.7                                 \
              /lib/x86_64-linux-ghc-8.10.7/libHSrts-1.0.1_thr-ghc8.10.7.so             \
              /lib/x86_64-linux-ghc-8.10.7/libHSstm-2.5.0.1-ghc8.10.7.so               \
              /lib/x86_64-linux-ghc-8.10.7/libHStemplate-haskell-2.16.0.0-ghc8.10.7.so \
              /lib/x86_64-linux-ghc-8.10.7/libHSterminfo-0.4.1.4-ghc8.10.7.so          \
              /lib/x86_64-linux-ghc-8.10.7/libHStime-1.9.3-ghc8.10.7.so                \
              /lib/x86_64-linux-ghc-8.10.7/libHStransformers-0.5.6.2-ghc8.10.7.so      \
              /lib/x86_64-linux-ghc-8.10.7/libHSunix-2.7.2.2-ghc8.10.7.so              /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/array-0.5.4.0                               /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/array-0.5.4.0
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/base-4.14.3.0                               /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/base-4.14.3.0
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/bytestring-0.10.12.0                        /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/bytestring-0.10.12.0
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/containers-0.6.5.1                          /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/containers-0.6.5.1
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/deepseq-1.4.4.0                             /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/deepseq-1.4.4.0
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/ghc-boot-th-8.10.7                          /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/ghc-boot-th-8.10.7
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/ghc-prim-0.6.1                              /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/ghc-prim-0.6.1
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/integer-simple-0.1.2.0                      /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/integer-simple-0.1.2.0
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/pretty-1.1.3.6                              /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/pretty-1.1.3.6
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/rts-1.0.1                                   /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/rts-1.0.1
COPY --from=0 /lib/x86_64-linux-ghc-8.10.7/template-haskell-2.16.0.0                   /usr/lib/ghc-8.10.7/lib/x86_64-linux-ghc-8.10.7/template-haskell-2.16.0.0
COPY --from=0 /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/cc1                     \
              /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/liblto_plugin.so        /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/local/bin/nibbles                                                   /usr/local/bin/
COPY --from=0 /usr/local/bin/Data/Digest/Murmur64.hi                                   \
              /usr/local/bin/Data/Digest/Murmur64.hs                                   \
              /usr/local/bin/Data/Digest/Murmur64.o                                    /usr/local/bin/Data/Digest/
COPY --from=0 /usr/local/bin/Data/Function/Memoize.hi                                  \
              /usr/local/bin/Data/Function/Memoize.hs                                  \
              /usr/local/bin/Data/Function/Memoize.o                                   /usr/local/bin/Data/Function/
COPY --from=0 /usr/local/bin/Data/Function/Memoize/Class.dyn_o                         \
              /usr/local/bin/Data/Function/Memoize/Class.hi                            \
              /usr/local/bin/Data/Function/Memoize/Class.hs                            \
              /usr/local/bin/Data/Function/Memoize/Class.o                             \
              /usr/local/bin/Data/Function/Memoize/TH.dyn_o                            \
              /usr/local/bin/Data/Function/Memoize/TH.hi                               \
              /usr/local/bin/Data/Function/Memoize/TH.hs                               \
              /usr/local/bin/Data/Function/Memoize/TH.o                                /usr/local/bin/Data/Function/Memoize/
COPY --from=0 /usr/local/bin/Data/List/Split.hi                                        \
              /usr/local/bin/Data/List/Split.hs                                        \
              /usr/local/bin/Data/List/Split.o                                         /usr/local/bin/Data/List/
COPY --from=0 /usr/local/bin/Data/List/Split/Internals.hi                              \
              /usr/local/bin/Data/List/Split/Internals.hs                              \
              /usr/local/bin/Data/List/Split/Internals.o                               /usr/local/bin/Data/List/Split/

ENTRYPOINT ["nibbleswrapper"]

CMD ["--version"]
