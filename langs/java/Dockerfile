# Also bump Rockstar.
FROM eclipse-temurin:23_37-jdk-alpine AS builder

RUN apk add --no-cache build-base gmp-dev

RUN jlink                                 \
    --add-modules java.base,jdk.compiler  \
    --compress zip-9                      \
    --module-path /opt/java/openjdk/jmods \
    --output /opt/jdk

COPY java.c /

RUN gcc -s -o java java.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /opt/jdk                 /opt/jdk
COPY --from=0 /java                    /usr/bin/

ENTRYPOINT ["java"]

CMD ["--version"]
