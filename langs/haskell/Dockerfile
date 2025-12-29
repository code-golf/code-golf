FROM alpine:3.23 AS builder

RUN apk add -X https://dl-cdn.alpinelinux.org/alpine/edge/community \
    --no-cache bash curl g++ ghc=9.10.3-r1 make ncurses-dev python3

ENV BOOTSTRAP_HASKELL_MINIMAL=1        \
    BOOTSTRAP_HASKELL_NONINTERACTIVE=1 \
    BOOTSTRAP_HASKELL_NO_UPGRADE=1     \
    PATH=/root/.ghcup/bin:$PATH        \
    # Also rebuild iogii.
    VER=9.14.1

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

WORKDIR /_build/stage1/lib

RUN find . \( -name '*.a'     \
          -or -name '*.hi'    \
          -or -name '*.p*hi'  \
          -or -name '*_p*.so' \
          -or -name '*debug*' \) -delete

COPY haskell.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s haskell.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/gcc              \
              /usr/bin/haskell          /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8      \
              /usr/lib/libgmp.so        \
              /usr/lib/libgmp.so.10     \
              /usr/lib/libncursesw.so.6 /usr/lib/
COPY --from=0 /_build/stage1/bin/ghc    /usr/local/bin/
COPY --from=0 /_build/stage1/lib        /usr/local/lib

ENTRYPOINT ["haskell"]

CMD ["--version"]
