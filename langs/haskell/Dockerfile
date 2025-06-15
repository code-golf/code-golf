ARG VER=9.12.2

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache bash build-base curl gmp-dev libffi-dev ncurses-dev python3

RUN curl https://get-ghcup.haskell.org \
  | BOOTSTRAP_HASKELL_GHC_VERSION=9.10 \
    BOOTSTRAP_HASKELL_NONINTERACTIVE=1 sh

ENV PATH=/root/.ghcup/bin:$PATH

RUN curl -#L https://downloads.haskell.org/~ghc/$VER/ghc-$VER-src.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure        \
 && hadrian/build      \
    -j`nproc`          \
    --docs=none        \
    --flavour=quickest \
    --freeze1

WORKDIR /_build/stage1

RUN strip bin/ghc \
          bin/runghc

# Delete static libraries.
RUN find lib -name '*.a' -delete

# Delete unused libraries.
RUN rm -r lib/x86_64-linux-ghc-$VER-inplace/Cabal-* \
          lib/x86_64-linux-ghc-$VER-inplace/ghci-*

COPY haskell.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s /haskell.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-x86_64.so.1                         /lib/
COPY --from=0 /usr/bin/gcc /usr/bin/haskell                    /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8                             \
              /usr/lib/libgmp.so                               \
              /usr/lib/libgmp.so.10                            \
              /usr/lib/libncursesw.so.6                        /usr/lib/
COPY --from=0 /_build/stage1/bin/ghc /_build/stage1/bin/runghc /usr/local/bin/
COPY --from=0 /_build/stage1/lib/settings                      /usr/local/lib/
COPY --from=0 /_build/stage1/lib/package.conf.d                /usr/local/lib/package.conf.d
COPY --from=0 /_build/stage1/lib/x86_64-linux-ghc-$VER-inplace /usr/local/lib/x86_64-linux-ghc-$VER-inplace

ENTRYPOINT ["haskell"]

CMD ["--version"]
