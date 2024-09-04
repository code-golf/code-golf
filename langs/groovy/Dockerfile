FROM ghcr.io/graalvm/graalvm-ce:22.3.1 AS builder

RUN microdnf install -y unzip

RUN curl -o groovy.zip https://groovy.jfrog.io/artifactory/dist-release-local/groovy-zips/apache-groovy-binary-4.0.22.zip

RUN unzip groovy.zip \
 && mv groovy-* /groovy

FROM codegolf/lang-base

ENV JAVA_HOME=/opt/jvm

COPY --from=0 /lib64       /lib64
COPY --from=0 /opt/graal*  /opt/jvm
COPY --from=0 /groovy/bin  \
              /usr/bin     /usr/bin/
COPY --from=0 /groovy/conf /usr/conf
COPY --from=0 /groovy/lib  /usr/lib

ENTRYPOINT ["groovy"]

CMD ["--version"]
