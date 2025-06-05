FROM alpine:3.22 AS builder

RUN apk add --no-cache bash curl

ENV VER=2.1.21

RUN curl -#Lo kotlin.zip https://github.com/JetBrains/kotlin/releases/download/v$VER/kotlin-compiler-$VER.zip

RUN unzip kotlin.zip && cp -r /kotlinc/* /usr/

# Remove some unneeded libraries.
RUN rm -f /usr/lib/android*.jar \
          /usr/lib/*sources.jar \
          /usr/lib/*test*.jar

FROM codegolf/lang-java

COPY --from=0 /bin/bash /bin/uname      /bin/
COPY --from=0 /usr/build.txt            /usr/
COPY --from=0 /usr/bin/dirname          \
              /usr/bin/env              \
              /usr/bin/kotlin           \
              /usr/bin/kotlinc          /usr/bin/
COPY --from=0 /usr/lib/*.jar            \
              /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/

FROM codegolf/lang-base

COPY --from=1 . /

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["kotlin"]

CMD ["-version"]
