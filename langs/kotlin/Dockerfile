FROM alpine:3.22 AS builder

RUN apk add --no-cache bash curl

ENV VER=2.2.21

RUN curl -#Lo kotlin.zip https://github.com/JetBrains/kotlin/releases/download/v$VER/kotlin-compiler-$VER.zip

RUN unzip kotlin.zip \
 && cp -r /kotlinc/* /usr

FROM codegolf/lang-java

COPY --from=0 /bin/bash /bin/uname                     /bin/
COPY --from=0 /usr/build.txt                           /usr/
COPY --from=0 /usr/bin/dirname                         \
              /usr/bin/env                             \
              /usr/bin/kotlin                          \
              /usr/bin/kotlinc                         /usr/bin/
COPY --from=0 /usr/lib/annotations-13.0.jar            \
              /usr/lib/kotlin-compiler.jar             \
              /usr/lib/kotlin-reflect.jar              \
              /usr/lib/kotlin-runner.jar               \
              /usr/lib/kotlin-script-runtime.jar       \
              /usr/lib/kotlin-scripting-*.jar          \
              /usr/lib/kotlin-stdlib.jar               \
              /usr/lib/kotlinx-coroutines-core-jvm.jar \
              /usr/lib/libncursesw.so.6                \
              /usr/lib/libreadline.so.8                /usr/lib/

SHELL ["/bin/bash", "-c"]

RUN /opt/java/bin/java -XX:+AutoCreateSharedArchive -XX:SharedArchiveFile=/cds.jsa \
    -jar /usr/lib/kotlin-compiler.jar -Xuse-fast-jar-file-system -e 'println(args)'

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["kotlin"]

CMD ["-version"]
