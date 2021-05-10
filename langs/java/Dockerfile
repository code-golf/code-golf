FROM adoptopenjdk/openjdk16:x86_64-debian-jdk-16.0.1_9-slim as builder

RUN mkdir /empty

RUN jlink                                 \
    --add-modules java.base,jdk.compiler  \
    --compress=2                          \
    --module-path /opt/java/openjdk/jmods \
    --output /opt/jdk

FROM alpine:3.13

RUN apk add --no-cache bash

FROM scratch

COPY --from=0 /bin/cat /bin/rm            /bin/
COPY --from=0 /lib/x86_64-linux-gnu       /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2 /lib64/
COPY --from=0 /opt/jdk                    /opt/jdk
COPY --from=0 /empty                      /proc
COPY --from=0 /empty                      /tmp
COPY --from=0 /usr/lib/locale             /usr/lib/locale/

# bash is used for the /usr/bin/java script. I couldn't get it to work with busybox sh.
COPY --from=1 /bin/bash                 /bin/
COPY --from=1 /lib/ld-musl-x86_64.so.1  /lib
COPY --from=1 /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/

COPY java /usr/bin/

ENTRYPOINT ["/opt/jdk/bin/java", "--version"]
