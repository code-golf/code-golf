ARG VER=5.0.3

FROM alpine:3.23 AS builder

ARG VER

RUN wget -qO groovy.zip https://groovy.jfrog.io/artifactory/dist-release-local/groovy-zips/apache-groovy-binary-$VER.zip \
 && unzip -j groovy

FROM codegolf/lang-java

ARG VER

COPY --from=0 /groovy              \
              /startGroovy         \
              /usr/bin/awk         \
              /usr/bin/basename    \
              /usr/bin/cut         \
              /usr/bin/expr        /usr/bin/
COPY --from=0 /groovy-starter.conf /usr/conf/
COPY --from=0 /groovy-$VER.jar     /usr/lib/

# Delete default CDS archive to save space.
RUN rm -f /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /bin /bin
COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["groovy"]

CMD ["--version"]
