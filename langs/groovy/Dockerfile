ARG VER=4.0.26

FROM eclipse-temurin:24_36-jre-alpine-3.21 AS builder

ARG VER

RUN wget -cO groovy.zip https://groovy.jfrog.io/artifactory/dist-release-local/groovy-zips/apache-groovy-binary-$VER.zip \
 && unzip -j groovy.zip

FROM codegolf/lang-base

ARG VER

COPY --from=0 /bin/sh /bin/uname                            /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
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
COPY --from=0 /groovy                                       \
              /startGroovy                                  \
              /usr/bin/awk                                  \
              /usr/bin/basename                             \
              /usr/bin/cut                                  \
              /usr/bin/dirname                              \
              /usr/bin/env                                  \
              /usr/bin/expr                                 /usr/bin/
COPY --from=0 /groovy-starter.conf                          /usr/conf/
COPY --from=0 /groovy-$VER.jar                              \
              /usr/lib/libacl.so.1                          \
              /usr/lib/libattr.so.1                         \
              /usr/lib/libcrypto.so.3                       \
              /usr/lib/libskarnet.so.2.14                   \
              /usr/lib/libutmps.so.0.1                      /usr/lib/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["groovy"]

CMD ["--version"]
