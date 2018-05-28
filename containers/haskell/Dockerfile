FROM alpine:edge

RUN apk add --no-cache ghc

RUN echo -e "#!/bin/sh -e\n\
\n\
/bin/cat - > /tmp/code.hs\n\
\n\
exec runghc /tmp/code.hs \$*" > /usr/bin/haskell

RUN chmod +x /usr/bin/haskell \
 && rm -rf /sbin /usr/share/perl5

ENTRYPOINT ["/usr/bin/ghc", "--version"]
