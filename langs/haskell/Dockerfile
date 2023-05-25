FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base ghc=9.4.4-r1

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

# Delete all the GHC binaries except ghc, runghc
RUN find /usr/lib/ghc-9.4.4/bin -type f \
    -not \( -name ghc-9.4.4 -or -name runghc-9.4.4 \) -delete

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /haskell                 \
              /usr/bin/gcc             /usr/bin/
COPY --from=0 /usr/lib                 /usr/lib

ENTRYPOINT ["haskell"]

CMD ["--version"]
