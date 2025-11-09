FROM alpine:3.22 AS builder

RUN apk add --no-cache coreutils

ENV VER=0.6.5

RUN wget -qO /usr/lib/cjam.jar https://downloads.sourceforge.net/cjam/cjam-$VER.jar

FROM codegolf/lang-java

COPY --from=0 /bin/sh                     /bin/
COPY --from=0 /usr/lib/cjam.jar           \
              /usr/lib/libacl.so.1        \
              /usr/lib/libattr.so.1       \
              /usr/lib/libcrypto.so.3     \
              /usr/lib/libskarnet.so.2.14 \
              /usr/lib/libutmps.so.0.1    /usr/lib/

# Delete default CDS archive to save space.
RUN rm /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr
