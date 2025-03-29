FROM eclipse-temurin:24_36-jre-alpine-3.21 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=0.6.5

RUN wget -O cjam.jar https://sourceforge.net/projects/cjam/files/cjam-$VER/cjam-$VER.jar/download

COPY cjam.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/cjam -s -static cjam.c

FROM codegolf/lang-base

COPY --from=0 /cjam.jar                              /
COPY --from=0 /lib/ld-musl-x86_64.so.1               /lib/
COPY --from=0 /opt/java/openjdk/bin/java             /opt/java/bin/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg          \
              /opt/java/openjdk/lib/libjava.so       \
              /opt/java/openjdk/lib/libjimage.so     \
              /opt/java/openjdk/lib/libjli.so        \
              /opt/java/openjdk/lib/libnet.so        \
              /opt/java/openjdk/lib/libnio.so        \
              /opt/java/openjdk/lib/libzip.so        \
              /opt/java/openjdk/lib/modules          /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so /opt/java/lib/server/
COPY --from=0 /usr/bin/cjam                          /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1                   \
              /usr/lib/libattr.so.1                  \
              /usr/lib/libcrypto.so.3                \
              /usr/lib/libskarnet.so.2.14            \
              /usr/lib/libutmps.so.0.1               /usr/lib/

ENTRYPOINT ["cjam"]

CMD ["--version"]
