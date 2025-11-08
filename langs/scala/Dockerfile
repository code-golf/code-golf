FROM alpine:3.22 AS builder

RUN apk add --no-cache bash coreutils curl

ENV VER=2.13.17

RUN curl -#L https://github.com/scala/scala/releases/download/v$VER/scala-$VER.tgz \
  | tar xz --directory /usr/local --strip-components 1

FROM codegolf/lang-java

COPY --from=0 /bin/bash /bin/grep /bin/uname    /bin/
COPY --from=0 /usr/bin/dirname /usr/bin/env     /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1              \
              /usr/lib/libattr.so.1             \
              /usr/lib/libcrypto.so.3           \
              /usr/lib/libncursesw.so.6         \
              /usr/lib/libreadline.so.8         \
              /usr/lib/libskarnet.so.2.14       \
              /usr/lib/libutmps.so.0.1          /usr/lib/
COPY --from=0 /usr/local/bin/scala              /usr/local/bin/
COPY --from=0 /usr/local/lib/scala-compiler.jar \
              /usr/local/lib/scala-library.jar  \
              /usr/local/lib/scala-reflect.jar  /usr/local/lib/

SHELL ["/bin/bash", "-c"]

# Delete default CDS archive used for AppCDS by Kotlin.
RUN rm /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /bin /bin
COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["scala"]

CMD ["--version"]
