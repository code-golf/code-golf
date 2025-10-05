FROM eclipse-temurin:25_36-jre-alpine-3.22 AS builder

RUN apk add --no-cache bash curl

ENV VER=2.2.20

RUN curl -#Lo kotlin.zip https://github.com/JetBrains/kotlin/releases/download/v$VER/kotlin-compiler-$VER.zip \
 && unzip kotlin

WORKDIR /kotlinc
# Remove some unneeded libraries.
RUN rm -f lib/android*.* \
          lib/*sources.* \
          lib/*test*.*   \
 && cp -r . /usr/

RUN java -XX:ArchiveClassesAtExit=/cds.jsa -jar /usr/lib/kotlin-compiler.jar -e ' '

FROM codegolf/lang-base

COPY --from=0 /cds.jsa                    /
COPY --from=0 /bin/bash /bin/uname        /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /opt/java/openjdk/bin/java  /opt/java/bin/
COPY --from=0 /opt/java/openjdk/lib       /opt/java/lib
COPY --from=0 /usr/build.txt              /usr/
COPY --from=0 /usr/bin/dirname            \
              /usr/bin/env                \
              /usr/bin/kotlin             \
              /usr/bin/kotlinc            /usr/bin/
COPY --from=0 /usr/lib/*.jar              \
              /usr/lib/libacl.so.1        \
              /usr/lib/libattr.so.1       \
              /usr/lib/libcrypto.so.3     \
              /usr/lib/libncursesw.so.6   \
              /usr/lib/libreadline.so.8   \
              /usr/lib/libskarnet.so.2.14 \
              /usr/lib/libutmps.so.0.1    /usr/lib/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["kotlin"]

CMD ["-version"]
