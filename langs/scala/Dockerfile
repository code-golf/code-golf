# For some reason it only finishes docker with AdoptOpenJDK 8 or 11
FROM adoptopenjdk/openjdk11:x86_64-alpine-jdk-11.28-slim AS builder

RUN apk add --no-cache curl

ENV PATH=/root/.local/share/coursier/bin:$PATH

RUN curl -fL https://github.com/coursier/coursier/releases/latest/download/cs-x86_64-pc-linux.gz | gzip -d > cs && chmod +x cs && ./cs setup

RUN curl -fL https://github.com/Virtuslab/scala-cli/releases/latest/download/scala-cli-x86_64-pc-linux.gz | gzip -d > scala-cli && chmod +x scala-cli && mv scala-cli /usr/bin/scala-cli

FROM codegolf/lang-base

COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /opt/java/openjdk                     /opt/java
COPY --from=0 /usr/bin/scala-cli                    /usr/bin/
COPY --from=0 /usr/glibc-compat/lib/libc.so.6       \
              /usr/glibc-compat/lib/libdl.so.2      \
              /usr/glibc-compat/lib/libpthread.so.0 \
              /usr/glibc-compat/lib/librt.so.1      \
              /usr/glibc-compat/lib/libz.so.1       /usr/glibc-compat/lib/

ENTRYPOINT ["scala-cli"]

CMD ["--version"]
