ARG VER=9.12.2

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache bash binutils-gold build-base cabal curl ghc-bootstrap gmp-dev ncurses-dev python3

RUN curl -#L https://downloads.haskell.org/~ghc/$VER/ghc-$VER-src.tar.xz \
  | tar xJ --strip-components 1

RUN cabal update                   \
 && ./configure                    \
 && hadrian/build -j               \
    --docs=none                    \
    --flavour=release+omit_pragmas \
    stage2:exe:ghc-bin

WORKDIR /_build/stage1/lib

RUN find . \( -name '*.a' -or -name '*debug*.so' -or -name '*p-ghc*.so' \) -delete

WORKDIR /_build/stage1/lib/x86_64-linux-ghc-$VER-inplace

RUN find . \( -name 'ghc-experimental*' -or -name 'ghci*' \) -exec rm -rf {} +

RUN find . -not \( -name '*Constants.h' -or -name '*.dyn_hi' -or -name '*ghc*.so' -or -name 'settings' \) -delete

COPY haskell.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s /haskell.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-*.so.1                              /lib/
COPY --from=0 /usr/bin/gcc /usr/bin/haskell                    /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8                             \
              /usr/lib/libgmp.so                               \
              /usr/lib/libgmp.so.10                            \
              /usr/lib/libncursesw.so.6                        /usr/lib/
COPY --from=0 /_build/stage1/bin/ghc                           /usr/local/bin/
COPY --from=0 /_build/stage1/lib/settings                      /usr/local/lib/
COPY --from=0 /_build/stage1/lib/package.conf.d                /usr/local/lib/package.conf.d
COPY --from=0 /_build/stage1/lib/x86_64-linux-ghc-$VER-inplace /usr/local/lib/x86_64-linux-ghc-$VER-inplace

ENTRYPOINT ["haskell"]

CMD ["--version"]
