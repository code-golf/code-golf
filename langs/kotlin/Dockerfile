FROM eclipse-temurin:23_37-jdk-alpine AS builder

RUN apk add --no-cache curl unzip bash build-base

RUN curl -L https://github.com/JetBrains/kotlin/releases/download/v2.0.21/kotlin-compiler-2.0.21.zip > kotlin.zip

RUN unzip kotlin.zip
RUN cp /kotlinc/lib/annotations-13.0.jar /kotlinc/bin/

COPY kotlin.c /
RUN gcc -o /usr/bin/kotlin -s -static kotlin.c

FROM codegolf/lang-base

COPY --from=0 /kotlinc/           /kotlinc/
COPY --from=0 /opt/java           /opt/jdk
COPY --from=0 /bin                /bin
COPY --from=0 /usr/               /usr/
COPY --from=0 /lib                /lib
COPY --from=0 /lib/ld-musl-*.so.1 /lib

ENTRYPOINT ["kotlin"]

CMD ["--version"]
