FROM ubuntu:25.04 AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl unzip

ENV VER=0.14.1

RUN curl -#Lo uiua.zip https://github.com/uiua-lang/uiua/releases/download/$VER/uiua-bin-x86_64-unknown-linux-gnu-no-audio.zip

RUN unzip uiua.zip -d /usr/bin

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6     \
              /lib/x86_64-linux-gnu/libgcc_s.so.1 \
              /lib/x86_64-linux-gnu/libm.so.6     /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2         /lib64/
COPY --from=0 /usr/bin/uiua                       /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
