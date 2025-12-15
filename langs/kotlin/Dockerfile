FROM alpine:3.23 AS builder

ENV VER=2.2.21

RUN wget -qO kotlin.zip https://github.com/JetBrains/kotlin/releases/download/v$VER/kotlin-compiler-$VER.zip

RUN unzip kotlin.zip \
 && cp -r /kotlinc/* /usr

FROM codegolf/lang-java

COPY --from=0 /usr/build.txt                           /usr/
COPY --from=0 /usr/bin/kotlin /usr/bin/kotlinc         /usr/bin/
COPY --from=0 /usr/lib/annotations-13.0.jar            \
              /usr/lib/kotlin-compiler.jar             \
              /usr/lib/kotlin-reflect.jar              \
              /usr/lib/kotlin-runner.jar               \
              /usr/lib/kotlin-script-runtime.jar       \
              /usr/lib/kotlin-scripting-*.jar          \
              /usr/lib/kotlin-stdlib.jar               \
              /usr/lib/kotlinx-coroutines-core-jvm.jar /usr/lib/

# Delete default CDS archive to save space.
RUN rm -f /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /bin /bin
COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["kotlin"]

CMD ["-version"]
