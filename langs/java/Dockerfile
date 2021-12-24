FROM eclipse-temurin:17.0.1_12-jdk-alpine as builder

RUN mkdir /empty

RUN apk add --no-cache build-base gmp-dev

RUN jlink                                 \
    --add-modules java.base,jdk.compiler  \
    --compress=2                          \
    --module-path /opt/java/openjdk/jmods \
    --output /opt/jdk

COPY java.c /

RUN gcc -s -o java java.c

FROM scratch

COPY --from=0 /java                    /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /lib/libz.so.1           /lib/
COPY --from=0 /opt/jdk                 /opt/jdk
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp

ENTRYPOINT ["/usr/bin/java"]

CMD ["--version"]
