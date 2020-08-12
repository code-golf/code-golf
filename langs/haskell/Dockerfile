FROM alpine:3.12 as builder

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/community --no-cache \
    ghc=8.8.4-r0

RUN echo -e "#!/bin/sh -e\n\
\n\
/bin/cat - > /tmp/code.hs\n\
\n\
exec runghc /tmp/code.hs \"\$@\"" > /usr/bin/haskell && chmod +x /usr/bin/haskell

# Slim down the image, this is all a bit hacky.
RUN find /usr/lib -name '*.a' -delete \
 && rm -rf                            \
    /home                             \
    /media                            \
    /mnt                              \
    /root                             \
    /run                              \
    /sbin                             \
    /srv                              \
    /usr/bin/dwp                      \
    /usr/bin/perl*                    \
    /usr/bin/pod*                     \
    /usr/lib/gcc                      \
    /usr/lib/ghc-8.8.4/bin/ghc-iserv* \
    /usr/lib/ghc-8.8.4/bin/haddock    \
    /usr/lib/ghc-8.8.4/Cabal*         \
    /usr/lib/libLLVM*                 \
    /usr/lib/perl5                    \
    /usr/libexec                      \
    /usr/share/perl5                  \
    /var

FROM scratch

COPY --from=0 / /

ENTRYPOINT ["/usr/bin/ghc", "--version"]
