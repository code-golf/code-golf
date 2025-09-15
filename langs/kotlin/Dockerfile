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

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["kotlin"]

CMD ["-version"]
