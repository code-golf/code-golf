FROM eclipse-temurin:23.0.1_11-jre-alpine AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=0.3.0.1

RUN wget -O kamilalisp.jar https://github.com/kspalaiologos/kamilalisp/releases/download/v$VER/kamilalisp-$VER.jar

COPY kamilalisp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/kamilalisp -s -static kamilalisp.c

FROM codegolf/lang-base

COPY --from=0 /kamilalisp.jar                               /
COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
COPY --from=0 /opt/java/openjdk/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg                 \
              /opt/java/openjdk/lib/libjava.so              \
              /opt/java/openjdk/lib/libjimage.so            \
              /opt/java/openjdk/lib/libjli.so               \
              /opt/java/openjdk/lib/libmanagement.so        \
              /opt/java/openjdk/lib/libmanagement_ext.so    \
              /opt/java/openjdk/lib/libnet.so               \
              /opt/java/openjdk/lib/libnio.so               \
              /opt/java/openjdk/lib/libzip.so               \
              /opt/java/openjdk/lib/modules                 /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so        /opt/java/lib/server/
COPY --from=0 /usr/bin/kamilalisp                           /usr/bin/

ENTRYPOINT ["kamilalisp"]

CMD ["--version"]
