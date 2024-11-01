ARG VER=4.0.23

FROM eclipse-temurin:23.0.1_11-jre-alpine AS builder

ARG VER

RUN apk add --no-cache curl

RUN curl -Lo groovy.zip https://groovy.jfrog.io/artifactory/dist-release-local/groovy-zips/apache-groovy-binary-$VER.zip

RUN unzip groovy.zip \
 && cp -r /groovy*/* /usr/

FROM codegolf/lang-base

ARG VER

COPY --from=0 /bin/sh /bin/uname                            /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                      \
              /lib/libacl.so.1                              \
              /lib/libattr.so.1                             \
              /lib/libcrypto.so.3                           \
              /lib/libskarnet.so.2.14                       \
              /lib/libutmps.so.0.1                          /lib/
COPY --from=0 /opt/java/openjdk/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg                 \
              /opt/java/openjdk/lib/libjava.so              \
              /opt/java/openjdk/lib/libjimage.so            \
              /opt/java/openjdk/lib/libjli.so               \
              /opt/java/openjdk/lib/libnet.so               \
              /opt/java/openjdk/lib/libnio.so               \
              /opt/java/openjdk/lib/libverify.so            \
              /opt/java/openjdk/lib/libzip.so               \
              /opt/java/openjdk/lib/modules                 \
              /opt/java/openjdk/lib/tzdb.dat                /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so        /opt/java/lib/server/
COPY --from=0 /usr/bin/awk                                  \
              /usr/bin/basename                             \
              /usr/bin/cut                                  \
              /usr/bin/dirname                              \
              /usr/bin/env                                  \
              /usr/bin/expr                                 \
              /usr/bin/groovy                               \
              /usr/bin/startGroovy                          /usr/bin/
COPY --from=0 /usr/conf/groovy-starter.conf                 /usr/conf/
COPY --from=0 /usr/lib/groovy-$VER.jar                      /usr/lib/

ENV JAVA_HOME=/opt/java

ENTRYPOINT ["groovy"]

CMD ["--version"]
