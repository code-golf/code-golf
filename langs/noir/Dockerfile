FROM rust:1.92.0-slim-trixie AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl git

ENV PATH=/root/.nargo/bin:$PATH VER=1.0.0

RUN curl -#L https://raw.githubusercontent.com/noir-lang/noirup/main/install | bash

RUN noirup --version v$VER-beta.17

RUN nargo new noir

WORKDIR /noir

RUN nargo check

COPY noir.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/noir -s noir.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/librt.so.1      /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --chown=nobody                                 \
     --from=0 /noir                                 /noir
COPY --from=0 /usr/bin/noir                         /usr/bin/
COPY --from=0 /root/.nargo/bin/nargo                /usr/local/bin/

ENTRYPOINT ["noir"]

CMD ["--version"]
