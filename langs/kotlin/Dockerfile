FROM eclipse-temurin:23.0.1_11-jre-alpine AS builder

RUN apk add --no-cache bash curl

ENV VER=2.0.21

RUN curl -L https://github.com/JetBrains/kotlin/releases/download/v$VER/kotlin-compiler-$VER.zip > kotlin.zip

RUN unzip kotlin.zip && cp -r /kotlinc/* /usr/

# Remove some unneeded libraries.
RUN rm -f /usr/lib/android*.jar \
          /usr/lib/*sources.jar \
          /usr/lib/*test*.jar

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/cat /bin/rm /bin/sh /bin/uname /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                      \
              /lib/libacl.so.1                              \
              /lib/libattr.so.1                             \
              /lib/libcrypto.so.3                           \
              /lib/libskarnet.so.2.14                       \
              /lib/libutmps.so.0.1                          \
              /lib/libz.so.1                                /lib/
COPY --from=0 /opt/java/openjdk/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jrt-fs.jar              \
              /opt/java/openjdk/lib/jvm.cfg                 \
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
COPY --from=0 /usr/build.txt                                /usr/
COPY --from=0 /usr/bin/dirname                              \
              /usr/bin/env                                  \
              /usr/bin/kotlin                               \
              /usr/bin/kotlinc                              /usr/bin/
COPY --from=0 /usr/lib/*.jar                                \
              /usr/lib/libncursesw.so.6                     \
              /usr/lib/libreadline.so.8                     /usr/lib/

COPY /kotlinwrapper /usr/bin/

ENV JAVA_HOME=/opt/java

ENTRYPOINT ["kotlinwrapper"]

CMD ["--version"]
