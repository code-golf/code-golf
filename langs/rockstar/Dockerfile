FROM eclipse-temurin:17.0.4.1_1-jre-alpine as builder

ENV VERSION=d23150cadc7e005ec2942f8d6e2646ca6f7e1f07

RUN apk add --no-cache wget

RUN wget -c https://github.com/gaborsch/rocky/archive/$VERSION.tar.gz -O - | tar -xz && \
    cd rocky-$VERSION && \
    (echo '#!/bin/sh' && echo 'exec /opt/java/openjdk/bin/java -jar $0 "$@"' && cat rocky.jar) > /usr/bin/rocky && \
    chmod +x /usr/bin/rocky

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/bin/rocky           \
              /usr/bin/dirname         \
              /usr/bin/basename        /usr/bin/
COPY --from=0 /bin/cat                 \
              /bin/rm                  \
              /bin/ln                  \
              /bin/sed                 \
              /bin/chmod               \
              /bin/sh                  /bin/
COPY --from=0 /opt/java/openjdk        /opt/java/openjdk

COPY rockstar /usr/bin

RUN chmod +x /usr/bin/rockstar

ENTRYPOINT ["rockstar"]

CMD ["-v"]
