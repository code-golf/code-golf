FROM ghcr.io/graalvm/graalvm-ce:22.3.1 AS builder

RUN curl -Lo scala-cli.gz https://github.com/VirtusLab/scala-cli/releases/download/v1.5.0/scala-cli-x86_64-pc-linux-static.gz

RUN gunzip scala-cli.gz    \
 && chmod +x scala-cli     \
 && mv scala-cli /usr/bin/ \
 && TERM="dumb" scala-cli  \
 && scala-cli --power bloop

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh   /bin/
COPY --from=0 /lib               /lib
COPY --from=0 /lib64             /lib64
COPY --from=0 /opt/graalvm-*     /opt/graalvm
COPY --from=0 /usr/bin/coreutils \
              /usr/bin/scala-cli /usr/bin/

COPY --from=0 /.cache            /.cache

COPY /scalawrapper /usr/bin/

# Remove
COPY --from=0 /bin/chmod /bin/
RUN chmod +x /usr/bin/scalawrapper

ENTRYPOINT ["scalawrapper"]

CMD ["--version"]
