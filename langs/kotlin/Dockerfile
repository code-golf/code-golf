FROM eclipse-temurin:23.0.1_11-jre-alpine AS builder

RUN apk add --no-cache bash curl unzip

ENV VER=2.0.21

RUN curl -L https://github.com/JetBrains/kotlin/releases/download/v$VER/kotlin-compiler-$VER.zip > kotlin.zip

RUN unzip kotlin.zip && cp -r /kotlinc/* /usr/

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/uname      /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1  \
              /lib/libacl.so.1          \
              /lib/libattr.so.1         \
              /lib/libcrypto.so.3       \
              /lib/libskarnet.so.2.14   \
              /lib/libutmps.so.0.1      \
              /lib/libz.so.1            /lib/
COPY --from=0 /opt/java                 /opt/java
COPY --from=0 /usr/build.txt            /usr/
COPY --from=0 /usr/bin/dirname          \
              /usr/bin/env              \
              /usr/bin/kotlin           \
              /usr/bin/kotlinc          /usr/bin/
COPY --from=0 /usr/lib/*.jar            \
              /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/

ENV PATH=$PATH:/opt/java/openjdk/bin

ENTRYPOINT ["kotlin"]

CMD ["-version"]
