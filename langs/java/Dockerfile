FROM adoptopenjdk/openjdk16:x86_64-debian-jdk-16.0.1_9-slim as builder

RUN mkdir /empty

RUN apt-get update && apt-get -y install gcc

RUN jlink                                 \
    --add-modules java.base,jdk.compiler  \
    --compress=2                          \
    --module-path /opt/java/openjdk/jmods \
    --output /opt/jdk

COPY java.c /

RUN gcc -s -o java java.c

FROM scratch

COPY --from=0 /java                       /usr/bin/
COPY --from=0 /lib/x86_64-linux-gnu       /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2 /lib64/
COPY --from=0 /opt/jdk                    /opt/jdk
COPY --from=0 /empty                      /proc
COPY --from=0 /empty                      /tmp
COPY --from=0 /usr/lib/locale             /usr/lib/locale/

ENTRYPOINT ["/usr/bin/java"]

CMD ["--version"]
