FROM eclipse-temurin:24.0.1_9-jre-alpine-3.21 AS builder

RUN apk add --no-cache bash curl gcc musl-dev

ENV VER=3.7.1

RUN curl -#L https://github.com/scala/scala3/releases/download/$VER/scala3-$VER.tar.gz \
  | tar xz --directory /usr/local --strip-components 1

COPY scala.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/scalawrapper -s -static scala.c

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/uname                          /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
COPY --from=0 /opt/java/openjdk/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jspawnhelper            \
              /opt/java/openjdk/lib/jvm.cfg                 \
              /opt/java/openjdk/lib/libjava.so              \
              /opt/java/openjdk/lib/libjimage.so            \
              /opt/java/openjdk/lib/libjli.so               \
              /opt/java/openjdk/lib/libmanagement.so        \
              /opt/java/openjdk/lib/libmanagement_ext.so    \
              /opt/java/openjdk/lib/libnet.so               \
              /opt/java/openjdk/lib/libnio.so               \
              /opt/java/openjdk/lib/libverify.so            \
              /opt/java/openjdk/lib/libzip.so               \
              /opt/java/openjdk/lib/modules                 \
              /opt/java/openjdk/lib/tzdb.dat                /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so        /opt/java/lib/server/
COPY --from=0 /usr/bin/dirname                              \
              /usr/bin/env                                  \
              /usr/bin/scalawrapper                         /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1                          \
              /usr/lib/libattr.so.1                         \
              /usr/lib/libcrypto.so.3                       \
              /usr/lib/libncursesw.so.6                     \
              /usr/lib/libreadline.so.8                     \
              /usr/lib/libskarnet.so.2.14                   \
              /usr/lib/libutmps.so.0.1                      /usr/lib/
COPY --from=0 /usr/local/VERSION                            /usr/local/
COPY --from=0 /usr/local/bin/scala                          /usr/local/bin/
COPY --from=0 /usr/local/libexec/cli-common-platform        \
              /usr/local/libexec/common-shared              \
              /usr/local/libexec/scala-cli.jar              /usr/local/libexec/
COPY --from=0 /usr/local/maven2/org                         /usr/local/maven2/org

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["scalawrapper"]

CMD ["--version"]
