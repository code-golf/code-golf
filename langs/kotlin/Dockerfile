FROM adoptopenjdk/openjdk11:x86_64-alpine-jdk-11.0.12_7-slim as builder

RUN mkdir /empty
RUN mkdir /myetc        \
 && echo nobody:x:65534:99::/: > /myetc/passwd

RUN apk add --no-cache curl zip unzip bash strace build-base linux-headers

RUN curl -L https://github.com/JetBrains/kotlin/releases/download/v2.0.21/kotlin-compiler-2.0.21.zip > kotlin.zip

RUN unzip kotlin.zip
RUN cp /kotlinc/lib/annotations-13.0.jar /kotlinc/bin/

COPY kotlin.c /
RUN gcc -o /usr/bin/kotlin -s -static kotlin.c

ENTRYPOINT ["/bin/bash"]

FROM scratch

COPY --from=0 /empty              /proc
COPY --from=0 /empty              /tmp
COPY --from=0 /empty              /dev
COPY --from=0 /myetc              /etc

COPY --from=0 /kotlinc/           /kotlinc/
COPY --from=0 /opt/java           /opt/jdk
COPY --from=0 /bin                /bin
COPY --from=0 /usr/               /usr/
COPY --from=0 /lib                /lib
COPY --from=0 /lib/ld-musl-*.so.1 /lib
COPY --from=0 /lib64              /lib64

ENTRYPOINT ["kotlin"]

CMD ["--version"]
