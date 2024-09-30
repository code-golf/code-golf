FROM alpine:3.20 AS builder

ENV VER=9.8.2

RUN apk add --no-cache gcc ghc=$VER-r1

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
RUN rm -r /usr/lib/ghc-$VER/lib/bin    \
 && find /usr/lib/ghc-$VER/bin -type f \
    -not \( -name ghc-$VER -or -name runghc-$VER \) -delete

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gcc             \
              /usr/bin/ghc             \
              /usr/bin/runghc          /usr/bin/
COPY --from=0 /usr/lib                 /usr/lib

COPY /haskell /usr/bin/

ENTRYPOINT ["haskell"]

CMD ["--version"]
