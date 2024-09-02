FROM ghcr.io/graalvm/graalvm-ce:22.3.1 AS builder
#FROM eclipse-temurin:22.0.2_9-jdk-alpine AS builder

#ENV VERSION=3.4.3

RUN microdnf install -y unzip
#RUN apk add --no-cache curl unzip

#RUN curl -Lo scala.zip https://github.com/scala/scala3/releases/download/$VERSION/scala3-$VERSION.zip

#RUN git clone https://github.com/Yewzir/code-golf-argv > argv

#RUN unzip scala.zip \
# && mv scala3-$VERSION /usr/local/scala

RUN curl -Lo scala.zip https://github.com/VirtusLab/scala-cli/releases/download/v1.5.0/scala-cli-x86_64-pc-linux-static-sdk.zip

RUN unzip scala.zip \
 && mv scala-* /usr/local/scala

FROM codegolf/lang-base

ENV JAVA_HOME=/opt/jvm

COPY --from=0 /bin             /bin
COPY --from=0 /lib             /lib
COPY --from=0 /lib64           /lib64
COPY --from=0 /opt/graalvm-*   /opt/jvm
COPY --from=0 /usr/bin         /usr/bin
COPY --from=0 /usr/lib         /usr/lib
COPY --from=0 /usr/local/scala /usr/local

COPY /scalawrapper /usr/bin/
RUN chmod +x /usr/bin/scalawrapper

ENTRYPOINT ["scalawrapper"]

CMD ["-version"]
