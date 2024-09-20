FROM ghcr.io/graalvm/graalvm-ce:22.3.1 AS builder

RUN curl -Lo scala-cli.gz https://github.com/VirtusLab/scala-cli/releases/download/v1.5.0/scala-cli-x86_64-pc-linux-static.gz

RUN gunzip scala-cli.gz    \
 && chmod +x scala-cli     \
 && mv scala-cli /usr/bin/ \
 && TERM="dumb" scala-cli  \
 && scala-cli --power bloop

FROM codegolf/lang-base

COPY --from=0 /bin           /bin
COPY --from=0 /lib           /lib
COPY --from=0 /lib64         /lib64
COPY --from=0 /opt/graalvm-* /opt/graalvm
COPY --from=0 /usr           /usr

COPY /scalawrapper /usr/bin/

# I know exactly what to do here but several attempts have failed so far.
# Well, "exactly" is too much, but I presume it'll work.

# Remove
RUN chmod +x /usr/bin/scalawrapper

ENTRYPOINT ["scala-cli"]

CMD ["--version"]
