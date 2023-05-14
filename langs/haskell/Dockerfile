FROM alpine:3.18 AS builder

RUN apk add --no-cache ghc=9.4.4-r1

RUN find /usr/lib -name '*.a' -delete

# Delete all libs except ffi, ghc, gmp, ncursesw
RUN find /usr/lib -not \(                    \
            -path '/usr/lib'                 \
        -or -path '/usr/lib/ghc-*'           \
        -or -path '/usr/lib/libffi.so*'      \
        -or -path '/usr/lib/libgmp.so*'      \
        -or -path '/usr/lib/libncursesw.so*' \
    \) -delete

# Delete all the /bin busybox symlinks except cat, sh
RUN find /bin -type l -not \( -name cat -or -name sh \) -delete

# Delete all the GHC binaries except ghc, runghc
RUN find /usr/lib/ghc-9.4.4/bin -type f \
    -not \( -name ghc-9.4.4 -or -name runghc-9.4.4 \) -delete

FROM codegolf/lang-base

COPY --from=0 /bin                     /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gcc             \
              /usr/bin/ghc             \
              /usr/bin/runghc          /usr/bin/
COPY --from=0 /usr/lib                 /usr/lib

COPY haskell /usr/bin/

ENTRYPOINT ["haskell"]

CMD ["--version"]
