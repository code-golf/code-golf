FROM alpine:edge

RUN apk add --no-cache ghc=8.4.3-r0

RUN echo -e "#!/bin/sh -e\n\
\n\
/bin/cat - > /tmp/code.hs\n\
\n\
exec runghc /tmp/code.hs \$*" > /usr/bin/haskell

# This is all a bit hacky
RUN chmod +x /usr/bin/haskell \
 && rm -rf                    \
    /home                     \
    /media                    \
    /mnt                      \
    /root                     \
    /run                      \
    /sbin                     \
    /srv                      \
    /usr/bin/perl*            \
    /usr/bin/pod*             \
    /usr/lib/gcc              \
    /usr/lib/perl5            \
    /usr/share/perl5          \
    /var

ENTRYPOINT ["/usr/bin/ghc", "--version"]
