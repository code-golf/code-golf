FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache ghc=8.8.4-r1

RUN find /usr/lib -name '*.a' -delete

# Delete some GHC stuff
RUN rm -r                     \
    /usr/bin/ghc-pkg*         \
    /usr/bin/ghci*            \
    /usr/lib/ghc-8.8.4/Cabal* \
    /usr/lib/ghc-8.8.4/html   \
    /usr/lib/ghc-8.8.4/latex

# Delete all libs except ffi, ghc, gmp, ncursesw
RUN find /usr/lib -not \(                    \
            -path '/usr/lib'                 \
        -or -path '/usr/lib/ghc-*'           \
        -or -path '/usr/lib/libffi.so*'      \
        -or -path '/usr/lib/libgmp.so*'      \
        -or -path '/usr/lib/libncursesw.so*' \
    \) -delete

# Delete all the /bin busybox symlinks except cat, sh
RUN find bin -type l -not \( -name cat -or -name sh \) -delete

# Delete all the GHC binaries except ghc, runghc
RUN find usr/lib/ghc-8.8.4/bin -type f -not \( -name ghc -or -name runghc \) -delete

# Delete all the /usr/bin files except gcc, ghc, runghc
RUN find usr/bin -not \( -name gcc -or -name 'ghc*' -or -name 'runghc*' \) -delete

FROM scratch

COPY --from=0 /bin                     /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin                 /usr/bin
COPY --from=0 /usr/lib                 /usr/lib

COPY haskell /usr/bin/

ENTRYPOINT ["/usr/bin/haskell"]

CMD ["--version"]
