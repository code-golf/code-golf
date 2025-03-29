FROM eclipse-temurin:24_36-jdk-alpine-3.21 AS builder

RUN apk add --no-cache bash curl git

ENV SBTVER=1.10.7 VER=2.13.16

RUN curl -#L https://github.com/sbt/sbt/releases/download/v$SBTVER/sbt-$SBTVER.tgz | tar xz

RUN git clone https://github.com/scala/scala.git --branch v$VER

WORKDIR /scala

RUN TERM='dumb' /sbt/bin/sbt dist/mkPack

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/grep /bin/uname                /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
COPY --from=0 /opt/java/openjdk/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg                 \
              /opt/java/openjdk/lib/libjava.so              \
              /opt/java/openjdk/lib/libjimage.so            \
              /opt/java/openjdk/lib/libjli.so               \
              /opt/java/openjdk/lib/libnet.so               \
              /opt/java/openjdk/lib/libnio.so               \
              /opt/java/openjdk/lib/libzip.so               \
              /opt/java/openjdk/lib/modules                 \
              /opt/java/openjdk/lib/tzdb.dat                /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so        /opt/java/lib/server/
COPY --from=0 /usr/bin/dirname /usr/bin/env                 /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1                          \
              /usr/lib/libattr.so.1                         \
              /usr/lib/libcrypto.so.3                       \
              /usr/lib/libgmp.so.10                         \
              /usr/lib/libncursesw.so.6                     \
              /usr/lib/libreadline.so.8                     \
              /usr/lib/libskarnet.so.2.14                   \
              /usr/lib/libutmps.so.0.1                      /usr/lib/
COPY --from=0 /scala/build/pack/bin/scala                   /usr/local/bin/
COPY --from=0 /scala/build/pack/lib/scala-compiler.jar      \
              /scala/build/pack/lib/scala-library.jar       \
              /scala/build/pack/lib/scala-reflect.jar       /usr/local/lib/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["scala"]

CMD ["--version"]
