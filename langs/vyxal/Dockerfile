FROM eclipse-temurin:24.0.1_9-jre-alpine-3.21 AS builder

ENV VER=3.8.0

RUN wget -O vyxal.jar https://github.com/Vyxal/Vyxal/releases/download/v$VER/vyxal-$VER.jar

FROM codegolf/lang-base

COPY --from=0 /vyxal.jar                             /
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

CMD ["/opt/java/bin/java", "-jar", "vyxal.jar", "-h"]
