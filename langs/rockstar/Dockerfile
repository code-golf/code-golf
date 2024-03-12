FROM eclipse-temurin:21.0.2_13-jre-alpine as builder

RUN apk add --no-cache curl

RUN curl -L https://github.com/gaborsch/rocky/archive/42e3a15.tar.gz  | tar xz \
 && cd rocky-*                                                                 \
 && ( echo '#!/bin/sh'                                                         \
      && echo 'exec /opt/java/openjdk/bin/java -jar $0 "$@"'                   \
      && cat rocky.jar                                                         \
    ) > /usr/bin/rocky                                                         \
 && chmod +x /usr/bin/rocky

FROM codegolf/lang-base

COPY --from=0 /bin/cat                   \
              /bin/sed                   \
              /bin/sh                    /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1   \
              /lib/libz.so.1             /lib/
COPY --from=0 /opt/java/openjdk/bin/java /opt/java/openjdk/bin/
COPY --from=0 /opt/java/openjdk/lib      /opt/java/openjdk/lib
COPY --from=0 /usr/bin/rocky             /usr/bin/

COPY rockstar /usr/bin

ENTRYPOINT ["rockstar"]

CMD ["-v"]
