ARG BUSYBOX_VER=1.37.0 VER=5.0.8

FROM alpine:3.24 AS builder

ARG BUSYBOX_VER VER

RUN apk add --no-cache coreutils build-base

RUN wget -qO groovy.zip https://groovy.jfrog.io/artifactory/dist-release-local/groovy-zips/apache-groovy-binary-$VER.zip \
 && unzip -j groovy

# Rewrite shebangs to bypass /usr/bin/env and go directly to /bin/sh.
RUN sed -i '1c#!/bin/sh' groovy startGroovy

# Groovy needs some UNIX utils, build BusyBox with symlinks to save space.
ADD --unpack https://busybox.net/downloads/busybox-$BUSYBOX_VER.tar.bz2 /

WORKDIR /busybox-$BUSYBOX_VER

COPY busybox-config .config

RUN make install

FROM codegolf/lang-java

ARG VER

COPY --from=0 /busybox-*/_install         /
COPY --from=0 /groovy                     \
              /startGroovy                /usr/bin/
COPY --from=0 /groovy-starter.conf        /usr/conf/
COPY --from=0 /groovy-$VER.jar            \
              /usr/lib/libacl.so.1        \
              /usr/lib/libattr.so.1       \
              /usr/lib/libcrypto.so.3     \
              /usr/lib/libskarnet.so.2.15 \
              /usr/lib/libutmps.so.0.1    /usr/lib/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["groovy"]

CMD ["--version"]
