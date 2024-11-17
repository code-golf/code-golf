FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl gmp-dev ncurses-dev

ENV VER=1.01

RUN curl https://golfscript.com/nibbles/nibbles-$VER.tgz \
  | tar xz --directory /usr/local/bin --strip-components 1

RUN curl -#L https://downloads.haskell.org/~ghcup/unofficial-bindists/ghc/8.10.7/ghc-8.10.7-x86_64-alpine-linux-integer-gmp.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure   \
    --prefix=/usr \
 && make install

WORKDIR /usr/local/bin

RUN ghc -O -package ghc *.hs

# IN PROGRESS: Strip whatever there's to strip here.

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/cp /bin/grep /bin/touch /bin/sh                     /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/libz.so.1                           /lib/
COPY --from=0 /usr/bin/as /usr/bin/gcc /usr/bin/ghc /usr/bin/ld                 /usr/bin/
COPY --from=0 /usr/include                                                      /usr/include
COPY --from=0 /usr/lib/crt1.o                                                   \
              /usr/lib/crti.o                                                   \
              /usr/lib/crtn.o                                                   \
              /usr/lib/libbfd-2.42.so                                           \
              /usr/lib/libc.so                                                  \
              /usr/lib/libctf.so.0                                              \
              /usr/lib/libdl.a                                                  \
              /usr/lib/libgcc_s.so                                              \
              /usr/lib/libgcc_s.so.1                                            \
              /usr/lib/libgmp.so                                                \
              /usr/lib/libgmp.so.10                                             \
              /usr/lib/libisl.so.23                                             \
              /usr/lib/libjansson.so.4                                          \
              /usr/lib/libm.a                                                   \
              /usr/lib/libmpc.so.3                                              \
              /usr/lib/libmpfr.so.6                                             \
              /usr/lib/libncursesw.so.6                                         \
              /usr/lib/librt.a                                                  \
              /usr/lib/libsframe.so.1                                           \
              /usr/lib/libssp_nonshared.a                                       \
              /usr/lib/libzstd.so.1                                             /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtbegin.o           \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtbeginS.o          \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtend.o             \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtendS.o            \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/
COPY --from=0 /usr/lib/ghc-8.10.7/platformConstants                             \
              /usr/lib/ghc-8.10.7/settings                                      /usr/lib/ghc-8.10.7/

# IN PROGRESS: Filter these directories!
COPY --from=0 /usr/lib/ghc-8.10.7/array-0.5.4.0                                 /usr/lib/ghc-8.10.7/array-0.5.4.0
COPY --from=0 /usr/lib/ghc-8.10.7/base-4.14.3.0                                 /usr/lib/ghc-8.10.7/base-4.14.3.0
COPY --from=0 /usr/lib/ghc-8.10.7/binary-0.8.8.0                                /usr/lib/ghc-8.10.7/binary-0.8.8.0
COPY --from=0 /usr/lib/ghc-8.10.7/bytestring-0.10.12.0                          /usr/lib/ghc-8.10.7/bytestring-0.10.12.0
COPY --from=0 /usr/lib/ghc-8.10.7/containers-0.6.5.1                            /usr/lib/ghc-8.10.7/containers-0.6.5.1
COPY --from=0 /usr/lib/ghc-8.10.7/deepseq-1.4.4.0                               /usr/lib/ghc-8.10.7/deepseq-1.4.4.0
COPY --from=0 /usr/lib/ghc-8.10.7/directory-1.3.6.0                             /usr/lib/ghc-8.10.7/directory-1.3.6.0
COPY --from=0 /usr/lib/ghc-8.10.7/exceptions-0.10.4                             /usr/lib/ghc-8.10.7/exceptions-0.10.4
COPY --from=0 /usr/lib/ghc-8.10.7/filepath-1.4.2.1                              /usr/lib/ghc-8.10.7/filepath-1.4.2.1
COPY --from=0 /usr/lib/ghc-8.10.7/ghc-8.10.7                                    /usr/lib/ghc-8.10.7/ghc-8.10.7
COPY --from=0 /usr/lib/ghc-8.10.7/ghc-boot-8.10.7                               /usr/lib/ghc-8.10.7/ghc-boot-8.10.7
COPY --from=0 /usr/lib/ghc-8.10.7/ghc-boot-th-8.10.7                            /usr/lib/ghc-8.10.7/ghc-boot-th-8.10.7
COPY --from=0 /usr/lib/ghc-8.10.7/ghc-heap-8.10.7                               /usr/lib/ghc-8.10.7/ghc-heap-8.10.7
COPY --from=0 /usr/lib/ghc-8.10.7/ghc-prim-0.6.1                                /usr/lib/ghc-8.10.7/ghc-prim-0.6.1
COPY --from=0 /usr/lib/ghc-8.10.7/ghci-8.10.7                                   /usr/lib/ghc-8.10.7/ghci-8.10.7
COPY --from=0 /usr/lib/ghc-8.10.7/haskeline-0.8.2                               /usr/lib/ghc-8.10.7/haskeline-0.8.2
COPY --from=0 /usr/lib/ghc-8.10.7/hpc-0.6.1.0                                   /usr/lib/ghc-8.10.7/hpc-0.6.1.0
COPY --from=0 /usr/lib/ghc-8.10.7/include                                       /usr/lib/ghc-8.10.7/include
COPY --from=0 /usr/lib/ghc-8.10.7/integer-gmp-1.0.3.0                           /usr/lib/ghc-8.10.7/integer-gmp-1.0.3.0
COPY --from=0 /usr/lib/ghc-8.10.7/mtl-2.2.2                                     /usr/lib/ghc-8.10.7/mtl-2.2.2
COPY --from=0 /usr/lib/ghc-8.10.7/package.conf.d                                /usr/lib/ghc-8.10.7/package.conf.d
COPY --from=0 /usr/lib/ghc-8.10.7/pretty-1.1.3.6                                /usr/lib/ghc-8.10.7/pretty-1.1.3.6
COPY --from=0 /usr/lib/ghc-8.10.7/process-1.6.13.2                              /usr/lib/ghc-8.10.7/process-1.6.13.2
COPY --from=0 /usr/lib/ghc-8.10.7/rts                                           /usr/lib/ghc-8.10.7/rts
COPY --from=0 /usr/lib/ghc-8.10.7/stm-2.5.0.1                                   /usr/lib/ghc-8.10.7/stm-2.5.0.1
COPY --from=0 /usr/lib/ghc-8.10.7/template-haskell-2.16.0.0                     /usr/lib/ghc-8.10.7/template-haskell-2.16.0.0
COPY --from=0 /usr/lib/ghc-8.10.7/terminfo-0.4.1.4                              /usr/lib/ghc-8.10.7/terminfo-0.4.1.4
COPY --from=0 /usr/lib/ghc-8.10.7/time-1.9.3                                    /usr/lib/ghc-8.10.7/time-1.9.3
COPY --from=0 /usr/lib/ghc-8.10.7/transformers-0.5.6.2                          /usr/lib/ghc-8.10.7/transformers-0.5.6.2
COPY --from=0 /usr/lib/ghc-8.10.7/unix-2.7.2.2                                  /usr/lib/ghc-8.10.7/unix-2.7.2.2

COPY --from=0 /usr/lib/ghc-8.10.7/bin/ghc                                       /usr/lib/ghc-8.10.7/bin/
COPY --from=0 /usr/libexec/gcc/x86_64-alpine-linux-musl/13.2.1/cc1              \
              /usr/libexec/gcc/x86_64-alpine-linux-musl/13.2.1/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/13.2.1/

# IN PROGRESS: Filter this directory.
COPY --from=0 /usr/local/bin                                                    /usr/local/bin

COPY nibbleswrapper /usr/bin/

ENTRYPOINT ["nibbleswrapper"]

CMD ["--version"]
