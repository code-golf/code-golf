FROM eclipse-temurin:24.0.1_9-jre-alpine-3.21 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=1.8.4

RUN curl -#Lo scala-cli.gz https://github.com/VirtusLab/scala-cli/releases/download/v$VER/scala-cli-x86_64-pc-linux-static.gz

RUN gunzip scala-cli.gz \
 && chmod +x scala-cli  \
 && mv scala-cli /usr/bin/

# Download dependencies to the local cache for offline use.
RUN TERM='dumb' scala-cli

COPY scala.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/scala -s -static scala.c

FROM codegolf/lang-base

COPY --from=0 /root/.cache/coursier/v1                      /cache
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
COPY --from=0 /usr/bin/scala /usr/bin/scala-cli             /usr/bin/

ENTRYPOINT ["scala"]

CMD ["--version"]
