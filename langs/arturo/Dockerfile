FROM debian:bookworm-slim AS builder

ENV PATH="/root/.nimble/bin:$PATH"

RUN apt-get update && apt-get install -y curl gcc git libpcre3-dev xz-utils

RUN curl https://nim-lang.org/choosenim/init.sh -sSf | sh -s -- -y && choosenim stable

RUN git clone https://github.com/arturo-lang/arturo.git

WORKDIR /arturo

RUN ./build.nims -i --log --mode mini

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6           \
              /lib/x86_64-linux-gnu/libdl.so.2          \
              /lib/x86_64-linux-gnu/libm.so.6           \
              /lib/x86_64-linux-gnu/libpcre.so.3        \
              /lib/x86_64-linux-gnu/libpthread.so.0     /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2               /lib64/
COPY --from=0 /arturo/bin/arturo                        /usr/bin/

ENTRYPOINT ["arturo"]

CMD ["-v"]
