ARG VER=3.6.1

FROM eclipse-temurin:23.0.1_11-jre-alpine AS builder

ARG VER

RUN apk add --no-cache curl

RUN wget -O- https://github.com/Vyxal/Vyxal/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./mill jvm.assembly

FROM codegolf/lang-base

ARG VER

COPY --from=0 /out/jvm/assembly.dest/vyxal-$VER.jar  /vyxal.jar
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
