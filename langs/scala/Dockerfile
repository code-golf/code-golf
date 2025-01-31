ARG VER=3.4.3

FROM eclipse-temurin:23.0.2_7-jre-alpine-3.21 AS builder

RUN apk add --no-cache bash curl gcc musl-dev

ARG VER

RUN curl -#L https://github.com/scala/scala3/releases/download/$VER/scala3-$VER.tar.gz \
  | tar xz --directory /usr/local --strip-components 1

COPY scala.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/scala -s -static scala.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /bin/bash /bin/uname                             /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                         /lib/
COPY --from=0 /opt/java/openjdk/bin/java                       /opt/java/bin/
COPY --from=0 /opt/java/openjdk/conf/security/java.security    /opt/java/conf/security/
COPY --from=0 /opt/java/openjdk/lib/jvm.cfg                    \
              /opt/java/openjdk/lib/libjava.so                 \
              /opt/java/openjdk/lib/libjimage.so               \
              /opt/java/openjdk/lib/libjli.so                  \
              /opt/java/openjdk/lib/libnet.so                  \
              /opt/java/openjdk/lib/libnio.so                  \
              /opt/java/openjdk/lib/libzip.so                  \
              /opt/java/openjdk/lib/modules                    /opt/java/lib/
COPY --from=0 /opt/java/openjdk/lib/server/libjvm.so           /opt/java/lib/server/
COPY --from=0 /usr/bin/dirname                                 \
              /usr/bin/env                                     \
              /usr/bin/scala                                   /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1                             \
              /usr/lib/libattr.so.1                            \
              /usr/lib/libcrypto.so.3                          \
              /usr/lib/libgmp.so.10                            \
              /usr/lib/libncursesw.so.6                        \
              /usr/lib/libreadline.so.8                        \
              /usr/lib/libskarnet.so.2.14                      \
              /usr/lib/libutmps.so.0.1                         /usr/lib/
COPY --from=0 /usr/local/bin/common /usr/local/bin/scala       /usr/local/bin/
COPY --from=0 /usr/local/lib/compiler-interface-1.9.6.jar      \
              /usr/local/lib/scala-asm-9.6.0-scala-1.jar       \
              /usr/local/lib/scala-library-2.13.12.jar         \
              /usr/local/lib/scala3-compiler_3-$VER.jar        \
              /usr/local/lib/scala3-interfaces-$VER.jar        \
              /usr/local/lib/scala3-library_3-$VER.jar         \
              /usr/local/lib/scala3-tasty-inspector_3-$VER.jar \
              /usr/local/lib/tasty-core_3-$VER.jar             /usr/local/lib/

ENV JAVA_HOME=/opt/java

ENTRYPOINT ["scala"]

CMD ["--version"]
