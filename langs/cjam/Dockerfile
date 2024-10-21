FROM eclipse-temurin:23_37-jre-alpine AS builder

ENV VER=0.6.5

RUN wget -O cjam.jar https://sourceforge.net/projects/cjam/files/cjam-$VER/cjam-$VER.jar/download

RUN echo $VER > version

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1               \
              /lib/libacl.so.1                       \
              /lib/libattr.so.1                      \
              /lib/libcrypto.so.3                    \
              /lib/libskarnet.so.2.14                \
              /lib/libutmps.so.0.1                   /lib/
COPY --from=0 /opt/java/openjdk/bin/java             /opt/jdk/bin/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg          \
              /opt/java/openjdk/lib/libjava.so       \
              /opt/java/openjdk/lib/libjimage.so     \
              /opt/java/openjdk/lib/libjli.so        \
              /opt/java/openjdk/lib/libnet.so        \
              /opt/java/openjdk/lib/libnio.so        \
              /opt/java/openjdk/lib/libzip.so        \
              /opt/java/openjdk/lib/modules          /opt/jdk/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so /opt/jdk/lib/server/
COPY --from=0 /cjam.jar /version                     /usr/bin/

COPY /cjam /usr/bin/

ENTRYPOINT ["cjam"]

CMD ["--version"]
