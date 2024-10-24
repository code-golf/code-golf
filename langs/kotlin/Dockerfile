FROM eclipse-temurin:23_37-jdk-alpine AS builder

RUN apk add --no-cache bash curl unzip

RUN curl -L https://github.com/JetBrains/kotlin/releases/download/v2.0.21/kotlin-compiler-2.0.21.zip > kotlin.zip

RUN unzip kotlin.zip
RUN cp /kotlinc/lib/annotations-13.0.jar /kotlinc/bin/

FROM codegolf/lang-base

COPY --from=0 /kotlinc/           /kotlinc/
COPY --from=0 /opt/java           /opt/jdk
COPY --from=0 /bin                /bin
COPY --from=0 /usr/               /usr/
COPY --from=0 /lib                /lib
COPY --from=0 /lib/ld-musl-*.so.1 /lib

ENV PATH=$PATH:/opt/jdk/openjdk/bin

ENTRYPOINT ["/kotlinc/bin/kotlin"]

CMD ["-version"]
