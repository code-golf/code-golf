# Also rebuild iogii.
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

RUN ./configure       \
 && hadrian/build     \
    -j`nproc`         \
    --docs=none       \
    --flavour=release \
    --freeze1

WORKDIR /_build/stage1/lib

# Delete all archive files.
RUN find . -name '*.a' -delete

# Delete useless shared object files.
RUN find . \( -name '*debug*.so' -or -name '*p-ghc*.so' \) -delete

WORKDIR /_build/stage1/lib/x86_64-linux-ghc-$VER-inplace

# Delete all files not matching a pattern.
RUN find . -type f -not \( -name '*Constants.h' -or -name '*.dyn_hi' -or -name '*ghc*.so' -or -name 'settings' \) -delete

# Delete unused package directories.
RUN rm -r Cabal* ghc-compact* ghc-experimental* ghc-toolchain* ghci* haddock* parsec* text* xhtml*

# Delete unused packages' shared object files, except GHCi's.
RUN rm *Cabal*.so *ghc-compact*.so *ghc-experimental*.so *ghc-toolchain*.so *haddock*.so *parsec*.so *text*.so *xhtml*.so

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
