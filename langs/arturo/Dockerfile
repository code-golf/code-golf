FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc xz-utils

ENV PATH=/root/.nimble/bin:$PATH VER=1feabea

RUN curl https://nim-lang.org/choosenim/init.sh | sh -s - -y \
 && choosenim stable

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

RUN ./build.nims --install --mode mini --release

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpcre.so.3    \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /root/.arturo/bin/arturo              /usr/bin/

ENTRYPOINT ["arturo"]

CMD ["--version"]
