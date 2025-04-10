FROM alpine:3.21 AS builder

# Also rebuild iogii.
ENV VER=9.8.2

RUN apk add --no-cache build-base ghc=$VER-r2

RUN echo "#define VER \"$VER\"" > haskell.h

COPY haskell.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haskell -s haskell.c

RUN find /usr/lib -name '*.a' -delete

# Delete all libraries except 'ffi', 'ghc', 'gmp', and 'ncursesw'.
RUN find /usr/lib -not \(                    \
            -path '/usr/lib'                 \
        -or -path '/usr/lib/ghc-*'           \
        -or -path '/usr/lib/libffi.so*'      \
        -or -path '/usr/lib/libgmp.so*'      \
        -or -path '/usr/lib/libncursesw.so*' \
    \) -delete

# Delete unused binaries and libraries.
RUN rm -r /usr/lib/ghc-$VER/lib/bin                           \
          /usr/lib/ghc-$VER/lib/x86_64-linux-ghc-$VER/Cabal-* \
          /usr/lib/ghc-$VER/lib/x86_64-linux-ghc-$VER/ghci-*  \
 && find /usr/lib/ghc-$VER/bin -type f                        \
    -not \( -name ghc-$VER -or -name runghc-$VER \) -delete

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gcc             \
              /usr/bin/haskell         /usr/bin/
COPY --from=0 /usr/lib                 /usr/lib

ENTRYPOINT ["haskell"]

CMD ["--version"]
