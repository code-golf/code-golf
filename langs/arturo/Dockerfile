FROM debian:bookworm-slim AS builder

ENV PATH="/root/.nimble/bin:$PATH" TERM=xterm

RUN apt-get update && apt-get install -y curl g++ gcc git gtk+-3.0 libgmp-dev libmpfr-dev webkit2gtk-4.0 xz-utils

RUN curl https://nim-lang.org/choosenim/init.sh -sSf | sh -s -- -y && choosenim stable

RUN git clone https://github.com/arturo-lang/arturo.git

WORKDIR /arturo

RUN ./build.nims --install --log

FROM codegolf/lang-base-no-proc

COPY --from=0 /root/.arturo/bin/arturo /usr/bin/

ENTRYPOINT ["arturo"]

CMD ["-v"]
