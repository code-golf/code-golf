# Also bump CJam, Groovy, Kotlin, Scala, Vyxal.
FROM eclipse-temurin:24_36-jdk-alpine-3.21 AS builder

RUN apk add --no-cache build-base gmp-dev

RUN jlink                                 \
    --add-modules java.base,jdk.compiler  \
    --compress zip-9                      \
    --module-path /opt/java/openjdk/jmods \
    --output /opt/jdk

COPY java.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/java -s -static java.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /opt/jdk                 /opt/jdk
COPY --from=0 /usr/bin/java            /usr/bin/
COPY --from=0 /usr/lib/libz.so.1       /usr/lib/

ENTRYPOINT ["java"]

CMD ["--version"]
