FROM alpine:3.23 AS builder

ENV VER=2.13.18

RUN wget -O- https://github.com/scala/scala/releases/download/v$VER/scala-$VER.tgz \
  | tar xz --directory /usr --strip-components 1

FROM codegolf/lang-java

COPY --from=0 /bin/grep                   /bin/
COPY --from=0 /usr/bin/scala              /usr/bin/
COPY --from=0 /usr/lib/scala-compiler.jar \
              /usr/lib/scala-library.jar  \
              /usr/lib/scala-reflect.jar  /usr/lib/

# Delete default CDS archive to save space.
RUN rm -f /opt/java/lib/server/classes.jsa

FROM codegolf/lang-base

COPY --from=1 /bin /bin
COPY --from=1 /lib /lib
COPY --from=1 /opt /opt
COPY --from=1 /usr /usr

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["scala"]

CMD ["--version"]
