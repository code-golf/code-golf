FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=1.9.0

RUN curl -#Lo scala-cli.gz https://github.com/VirtusLab/scala-cli/releases/download/v$VER/scala-cli-x86_64-pc-linux-static.gz

RUN gunzip scala-cli.gz \
 && chmod +x scala-cli  \
 && mv scala-cli /usr/bin/

# Download dependencies to the local cache for offline use.
RUN TERM='dumb' scala-cli

# Delete some previously downloaded dependencies.
RUN find ~/.cache -type d -name 'github.com' -exec rm -r {} +

COPY scala.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/scala -s -static scala.c

FROM codegolf/lang-java

COPY --from=0 /root/.cache/coursier/v1 /cache
COPY --from=0 /usr/bin/scala           \
              /usr/bin/scala-cli       /usr/bin/

ENV JAVA_HOME='/opt/java'

ENTRYPOINT ["scalac"]

CMD ["--version"]
