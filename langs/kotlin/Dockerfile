FROM codegolf/lang-java

FROM alpine:3.22 AS builder

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

COPY --from=0 /opt /opt

RUN /opt/java/bin/java -XX:ArchiveClassesAtExit=/cds.jsa -jar /usr/lib/kotlin-compiler.jar -Xuse-fast-jar-file-system -e 'println(args)'

FROM codegolf/lang-java

COPY --from=1 /cds.jsa                  /
COPY --from=1 /bin/bash /bin/uname      /bin/
COPY --from=1 /usr/build.txt            /usr/
COPY --from=1 /usr/bin/dirname          \
              /usr/bin/env              \
              /usr/bin/kotlin           \
              /usr/bin/kotlinc          /usr/bin/
COPY --from=1 /usr/lib/*.jar            \
              /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["kotlin"]

CMD ["-version"]
