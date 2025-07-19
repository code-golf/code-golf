ARG VER=4.0.27

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache coreutils dash

RUN wget -qO groovy.zip https://groovy.jfrog.io/artifactory/dist-release-local/groovy-zips/apache-groovy-binary-$VER.zip \
 && unzip -j groovy

FROM codegolf/lang-java

ARG VER

COPY --from=0 /bin/uname                  /bin/
COPY --from=0 /usr/bin/dash               /bin/sh
COPY --from=0 /groovy                     \
              /startGroovy                \
              /usr/bin/awk                \
              /usr/bin/basename           \
              /usr/bin/cut                \
              /usr/bin/dirname            \
              /usr/bin/env                \
              /usr/bin/expr               /usr/bin/
COPY --from=0 /groovy-starter.conf        /usr/conf/
COPY --from=0 /groovy-$VER.jar            \
              /usr/lib/libacl.so.1        \
              /usr/lib/libattr.so.1       \
              /usr/lib/libcrypto.so.3     \
              /usr/lib/libskarnet.so.2.14 \
              /usr/lib/libutmps.so.0.1    /usr/lib/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["groovy"]

CMD ["--version"]
