FROM alpine:3.20 AS builder

RUN apk add --no-cache bash curl git openjdk21-jre-headless

RUN curl -L https://github.com/sbt/sbt/releases/download/v1.10.1/sbt-1.10.1.tgz | tar xz \
 && mv sbt /usr/local/sbt

ENV PATH=/usr/local/sbt/bin:$PATH

RUN git clone https://github.com/scala/scala3.git -b 3.4.2 \
 && mv scala3 /usr/local/scala

WORKDIR /usr/local/scala
RUN sbt dist/packArchive

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/busybox /bin/cat /bin/chmod /bin/sh /bin/uname        /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/libz.so.1                              /lib/
COPY          /scalarunner                                                         /usr/bin/
COPY --from=0 /usr/bin/dirname /usr/bin/env                                        /usr/bin/
COPY --from=0 /usr/lib/jvm/java-21-openjdk/jre/bin/java                            /usr/lib/jre/bin/
COPY --from=0 /usr/lib/jvm/java-21-openjdk/jre/conf/security/java.security         /usr/lib/jre/conf/security/
COPY --from=0 /usr/lib/jvm/java-21-openjdk/jre/lib/jvm.cfg                         \
              /usr/lib/jvm/java-21-openjdk/jre/lib/libjava.so                      \
              /usr/lib/jvm/java-21-openjdk/jre/lib/libjimage.so                    \
              /usr/lib/jvm/java-21-openjdk/jre/lib/libjli.so                       \
              /usr/lib/jvm/java-21-openjdk/jre/lib/libnet.so                       \
              /usr/lib/jvm/java-21-openjdk/jre/lib/libnio.so                       \
              /usr/lib/jvm/java-21-openjdk/jre/lib/libzip.so                       \
              /usr/lib/jvm/java-21-openjdk/jre/lib/modules                         /usr/lib/jre/lib/
COPY --from=0 /usr/lib/jvm/java-21-openjdk/jre/lib/server/libjvm.so                /usr/lib/jre/lib/server/
COPY --from=0 /usr/lib/libncursesw.so.6 /usr/lib/libreadline.so.8                  /usr/lib/
COPY --from=0 /usr/local/scala/dist/target/pack/bin/common                         \
              /usr/local/scala/dist/target/pack/bin/scala                          \
              /usr/local/scala/dist/target/pack/bin/scalac                         /usr/local/bin/
COPY --from=0 /usr/local/scala/dist/target/pack/lib/compiler-interface-1.9.6.jar   \
              /usr/local/scala/dist/target/pack/lib/scala-asm-9.6.0-scala-1.jar    \
              /usr/local/scala/dist/target/pack/lib/scala-library-2.13.12.jar      \
              /usr/local/scala/dist/target/pack/lib/scala3-compiler_3-*.jar        \
              /usr/local/scala/dist/target/pack/lib/scala3-interfaces-*.jar        \
              /usr/local/scala/dist/target/pack/lib/scala3-library_3-*.jar         \
              /usr/local/scala/dist/target/pack/lib/scala3-tasty-inspector_3-*.jar \
              /usr/local/scala/dist/target/pack/lib/tasty-core_3-*.jar             /usr/local/lib/

RUN chmod +x /usr/bin/scalarunner

ENTRYPOINT ["scalarunner"]

CMD ["-version"]
