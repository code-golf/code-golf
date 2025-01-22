FROM eclipse-temurin:23.0.2_7-jre-alpine-3.21 AS builder

RUN apk add --no-cache build-base wget

RUN wget https://github.com/gaborsch/rocky/raw/9af6921/rocky.jar

COPY rockstar.c .

RUN gcc -Wall -Werror -Wextra -o rockstar -s rockstar.c

FROM codegolf/lang-base

COPY --from=0 /rocky.jar                                    /
COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
COPY --from=0 /opt/java/openjdk/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg                 \
              /opt/java/openjdk/lib/libjava.so              \
              /opt/java/openjdk/lib/libjimage.so            \
              /opt/java/openjdk/lib/libjli.so               \
              /opt/java/openjdk/lib/libnet.so               \
              /opt/java/openjdk/lib/libnio.so               \
              /opt/java/openjdk/lib/libverify.so            \
              /opt/java/openjdk/lib/libzip.so               \
              /opt/java/openjdk/lib/modules                 /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so        /opt/java/lib/server/
COPY --from=0 /rockstar                                     /usr/bin/

ENTRYPOINT ["rockstar"]

CMD ["-h"]
