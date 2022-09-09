FROM debian:bullseye-slim as builder

RUN apt-get update && apt-get install -y curl libarchive-tools

RUN curl -L https://github.com/wren-lang/wren-cli/releases/download/0.4.0/wren-cli-linux-0.4.0.zip | bsdtar -xf -

RUN chmod +x /wren-cli-linux-0.4.0/wren_cli

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                      /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /wren-cli-linux-0.4.0/wren_cli        /usr/bin/

COPY wren /usr/bin/

ENTRYPOINT ["wren"]

CMD ["--version"]
