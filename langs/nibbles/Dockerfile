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

COPY --from=0 /bin/cat /bin/cp /bin/touch /bin/sh   /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1              \
              /lib/libz.so.1                        /lib/
COPY --from=0 /usr/bin/as                           \
              /usr/bin/gcc                          \
              /usr/bin/ghc                          \
              /usr/bin/ld                           /usr/bin/
COPY --from=0 /usr/include                          /usr/include
COPY --from=0 /usr/lib/crt1.o                       \
              /usr/lib/crti.o                       \
              /usr/lib/crtn.o                       \
              /usr/lib/libbfd-2.42.so               \
              /usr/lib/libc.so                      \
              /usr/lib/libctf.so.0                  \
              /usr/lib/libdl.a                      \
              /usr/lib/libgcc_s.so                  \
              /usr/lib/libgcc_s.so.1                \
              /usr/lib/libgmp.so                    \
              /usr/lib/libgmp.so.10                 \
              /usr/lib/libisl.so.23                 \
              /usr/lib/libjansson.so.4              \
              /usr/lib/libm.a                       \
              /usr/lib/libmpc.so.3                  \
              /usr/lib/libmpfr.so.6                 \
              /usr/lib/libncursesw.so.6             \
              /usr/lib/librt.a                      \
              /usr/lib/libsframe.so.1               \
              /usr/lib/libssp_nonshared.a           \
              /usr/lib/libzstd.so.1                 \
              /usr/lib/gcc/*/*/crtbegin.o           \
              /usr/lib/gcc/*/*/crtbeginS.o          \
              /usr/lib/gcc/*/*/crtend.o             \
              /usr/lib/gcc/*/*/crtendS.o            \
              /usr/lib/gcc/*/*/libgcc.a             /usr/lib/

# IN PROGRESS: Filtering this gigantic directory, likely à la Haskell.
COPY --from=0 /usr/lib/ghc-8.10.7                   /usr/lib/ghc-8.10.7

COPY --from=0 /usr/libexec/gcc/*/*/cc1              \
              /usr/libexec/gcc/*/*/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/13.2.1/

# IN PROGRESS: Filtering this directory. I can't stand integral directory copies like /usr/include really.
COPY --from=0 /usr/local/bin                        /usr/local/bin

COPY nibbleswrapper /usr/bin/

ENTRYPOINT ["nibbleswrapper"]

CMD ["--version"]
