FROM alpine:3.20 AS builder

RUN apk add --no-cache bash curl openjdk8-jre

RUN curl -L https://github.com/scala/scala3/releases/download/3.4.2/scala3-3.4.2.tar.gz | tar xz --strip-components=1

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/busybox /bin/common    \
              /bin/scala /bin/scalac /bin/uname     /bin/
COPY --from=0 /lib/compiler-interface-1.9.6.jar     \                            
              /lib/ld-musl-x86_64.so.1              \                  
              /lib/libz.so.1                        \
              /lib/scala-library-2.13.12.jar        \
              /lib/scala3-compiler_3-3.4.2.jar      \
              /lib/scala3-interfaces-3.4.2.jar      \
              /lib/scala3-library_3-3.4.2.jar       /lib/
COPY --from=0 /usr/bin/dirname /usr/bin/env         /usr/bin/  
COPY --from=0 /usr/lib/jvm/java-1.8-openjdk/jre/bin /usr/lib/jvm/java-1.8-openjdk/jre/bin
COPY --from=0 /usr/lib/jvm/java-1.8-openjdk/jre/lib /usr/lib/jvm/java-1.8-openjdk/jre/lib
COPY --from=0 /usr/lib/libgcc_s.so.1                \
              /usr/lib/libncursesw.so.6             \
              /usr/lib/libreadline.so.8             \
              /usr/lib/libstdc++.so.6               /usr/lib/

ENV JAVA_HOME="/usr/lib/jvm/java-1.8-openjdk/jre"

ENTRYPOINT ["scala"]

CMD ["-version"]
