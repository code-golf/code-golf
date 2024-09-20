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

COPY /scalawrapper /usr/bin/

# Remove 23-24
COPY --from=0 /bin/chmod /bin/
RUN chmod +x /usr/bin/scalawrapper

RUN chmod 777 /tmp
USER nobody
COPY --chown=nobody --from=0 /.cache /.cache

ENTRYPOINT ["scalawrapper"]

CMD ["--version"]
