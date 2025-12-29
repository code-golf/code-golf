ARG VER=9.14.1

FROM alpine:3.23 AS builder

ARG VER

RUN apk add -X https://dl-cdn.alpinelinux.org/alpine/edge/community \
    --no-cache bash curl g++ ghc=9.10.3-r1 make ncurses-dev python3

ENV BOOTSTRAP_HASKELL_MINIMAL=1        \
    BOOTSTRAP_HASKELL_NONINTERACTIVE=1 \
    BOOTSTRAP_HASKELL_NO_UPGRADE=1     \
    PATH=/root/.ghcup/bin:$PATH

RUN curl https://get-ghcup.haskell.org | sh

RUN curl -#L https://downloads.haskell.org/~ghc/$VER/ghc-$VER-src.tar.xz \
  | tar xJ --strip-components 1

RUN ghcup install cabal \
 && cabal update        \
 && ./configure         \
 && hadrian/build       \
    --docs=none         \
    --flavour=release   \
    --jobs=`nproc`      \
    stage2:exe:ghc-bin

WORKDIR /_build/stage1/lib/x86_64-linux-ghc-$VER-d531

RUN find . \( -name '*.a'     \
          -or -name '*.hi'    \
          -or -name '*.p*hi'  \
          -or -name '*_p*.so' \
          -or -name '*debug*' \) -delete

RUN strip *.so

COPY haskell.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s haskell.c

FROM codegolf/lang-base

ARG VER

ENV GHC=lib/x86_64-linux-ghc-$VER-d531

COPY --from=0 /lib/ld-musl-*.so.1                      /lib/
COPY --from=0 /usr/bin/gcc /usr/bin/haskell            /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8                     \
              /usr/lib/libgmp.so                       \
              /usr/lib/libgmp.so.10                    \
              /usr/lib/libncursesw.so.6                /usr/lib/
COPY --from=0 /_build/stage1/bin/ghc                   /usr/local/bin/
COPY --from=0 /_build/stage1/lib/settings              /usr/local/lib/
COPY --from=0 /_build/stage1/lib/package.conf.d        /usr/local/lib/package.conf.d
COPY --from=0 /_build/stage1/$GHC/libHS*-*_thr-*.so    \
              /_build/stage1/$GHC/libHS*-*-*.so        \
              /_build/stage1/$GHC/libffi*.so           /usr/local/$GHC/
COPY --from=0 /_build/stage1/$GHC/array-*-*            /usr/local/$GHC/array-0.5.8.0-d432
COPY --from=0 /_build/stage1/$GHC/base-*-*             /usr/local/$GHC/base-4.22.0.0-d466
COPY --from=0 /_build/stage1/$GHC/binary-*-*           /usr/local/$GHC/binary-0.8.9.3-0c31
COPY --from=0 /_build/stage1/$GHC/bytestring-*-*       /usr/local/$GHC/bytestring-0.12.2.0-8899
COPY --from=0 /_build/stage1/$GHC/containers-*-*       /usr/local/$GHC/containers-0.8-3e81
COPY --from=0 /_build/stage1/$GHC/deepseq-*-*          /usr/local/$GHC/deepseq-1.5.1.0-dd3e
COPY --from=0 /_build/stage1/$GHC/directory-*-*        /usr/local/$GHC/directory-1.3.10.0-0e26
COPY --from=0 /_build/stage1/$GHC/exceptions-*-*       /usr/local/$GHC/exceptions-0.10.11-82ea
COPY --from=0 /_build/stage1/$GHC/file-io-*-*          /usr/local/$GHC/file-io-0.1.5-e8d1
COPY --from=0 /_build/stage1/$GHC/filepath-*-*         /usr/local/$GHC/filepath-1.5.4.0-9b98
COPY --from=0 /_build/stage1/$GHC/ghc-*-*              /usr/local/$GHC/ghc-$VER-c343
COPY --from=0 /_build/stage1/$GHC/ghc-boot-*-*         /usr/local/$GHC/ghc-boot-$VER-dfe8
COPY --from=0 /_build/stage1/$GHC/ghc-boot-th-*-*      /usr/local/$GHC/ghc-boot-th-$VER-9249
COPY --from=0 /_build/stage1/$GHC/ghc-heap-*-*         /usr/local/$GHC/ghc-heap-$VER-c4eb
COPY --from=0 /_build/stage1/$GHC/ghc-internal-*-*     /usr/local/$GHC/ghc-internal-9.1401.0-ec29
COPY --from=0 /_build/stage1/$GHC/ghc-platform-*-*     /usr/local/$GHC/ghc-platform-0.1.0.0-e893
COPY --from=0 /_build/stage1/$GHC/ghc-prim-*-*         /usr/local/$GHC/ghc-prim-0.13.1-ee6a
COPY --from=0 /_build/stage1/$GHC/ghci-*-*             /usr/local/$GHC/ghci-$VER-b0b2
COPY --from=0 /_build/stage1/$GHC/hpc-*-*              /usr/local/$GHC/hpc-0.7.0.2-d87b
COPY --from=0 /_build/stage1/$GHC/mtl-*-*              /usr/local/$GHC/mtl-2.3.1-d5b6
COPY --from=0 /_build/stage1/$GHC/os-string-*-*        /usr/local/$GHC/os-string-2.0.8-3116
COPY --from=0 /_build/stage1/$GHC/pretty-*-*           /usr/local/$GHC/pretty-1.1.3.6-0717
COPY --from=0 /_build/stage1/$GHC/process-*-*          /usr/local/$GHC/process-1.6.26.1-f719
COPY --from=0 /_build/stage1/$GHC/rts-*                /usr/local/$GHC/rts-1.0.3
COPY --from=0 /_build/stage1/$GHC/semaphore-compat-*-* /usr/local/$GHC/semaphore-compat-1.0.0-5b9a
COPY --from=0 /_build/stage1/$GHC/stm-*-*              /usr/local/$GHC/stm-2.5.3.1-d5a1
COPY --from=0 /_build/stage1/$GHC/template-haskell-*-* /usr/local/$GHC/template-haskell-2.24.0.0-0225
COPY --from=0 /_build/stage1/$GHC/time-*-*             /usr/local/$GHC/time-1.15-1bfa
COPY --from=0 /_build/stage1/$GHC/transformers-*-*     /usr/local/$GHC/transformers-0.6.1.2-bcd0
COPY --from=0 /_build/stage1/$GHC/unix-*-*             /usr/local/$GHC/unix-2.8.8.0-11df

ENTRYPOINT ["haskell"]

CMD ["--version"]
