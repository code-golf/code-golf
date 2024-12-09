FROM alpine:3.21 AS builder

ENV VER=9.8.2

RUN apk add --no-cache build-base ghc=$VER-r2

RUN echo "#define VERSION \"$VER\"" > haskell.h

COPY haskell.c .

RUN gcc -Wall -Werror -Wextra -o haskell -s haskell.c

RUN find /usr/lib -name '*.a' -delete

# Delete all libs except ffi, ghc, gmp, ncursesw
RUN find /usr/lib -not \(                    \
            -path '/usr/lib'                 \
        -or -path '/usr/lib/ghc-*'           \
        -or -path '/usr/lib/libffi.so*'      \
        -or -path '/usr/lib/libgmp.so*'      \
        -or -path '/usr/lib/libncursesw.so*' \
    \) -delete

# Delete unused libs and binaries.
RUN rm -r /usr/lib/ghc-$VER/lib/bin                           \
          /usr/lib/ghc-$VER/lib/x86_64-linux-ghc-$VER/Cabal-* \
          /usr/lib/ghc-$VER/lib/x86_64-linux-ghc-$VER/ghci-*  \
 && find /usr/lib/ghc-$VER/bin -type f                        \
    -not \( -name ghc-$VER -or -name runghc-$VER \) -delete

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /haskell                 \
              /usr/bin/gcc             /usr/bin/
COPY --from=0 /usr/lib                 /usr/lib

ENTRYPOINT ["haskell"]

CMD ["--version"]
