FROM alpine:3.20 AS builder

ENV VER=1.01

RUN apk add --no-cache build-base curl gmp-dev ncurses-dev

ENV BOOTSTRAP_HASKELL_GHC_VERSION=8.10.7 \
    BOOTSTRAP_HASKELL_NONINTERACTIVE=1   \
    GHCUP_INSTALL_BASE_PREFIX=/usr       \
    PATH=/usr/.ghcup/bin:$PATH

RUN curl https://get-ghcup.haskell.org | sh

RUN curl https://golfscript.com/nibbles/nibbles-$VER.tgz \
  | tar xz --directory /usr/local/bin --strip-components 1

WORKDIR /usr/local/bin

RUN ghc -O -package ghc *.hs

FROM codegolf/lang-base

COPY --from=0 /bin                      /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/.ghcup/bin/runghc    /usr/bin/
COPY --from=0 /usr/lib/libgmp.so.10     \
              /usr/lib/libncursesw.so.6 /usr/lib/
COPY --from=0 /usr/local/bin            /usr/local/bin

COPY nibbleswrapper /usr/bin/

ENTRYPOINT ["nibbleswrapper"]

CMD ["--version"]
