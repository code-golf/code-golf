FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl g++ git make unzip

ENV VER=0.100

RUN git clone --depth 1 --branch $VER --progress https://github.com/factor/factor

RUN curl -#o factor.zip https://downloads.factorcode.org/releases/$VER/factor-src-$VER.zip \
 && unzip -n factor

WORKDIR /factor

RUN sed -i '/^\[/,/al/d' basis/vocabs/prettyprint/prettyprint.factor

RUN ./build.sh compile \
 && ./build.sh bootstrap

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/libm.so.6      \
              /lib/x86_64-linux-gnu/libstdc++.so.6 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=0 /factor/factor /factor/factor.image  /usr/bin/
COPY --from=0 /factor/basis                        /usr/bin/basis
COPY --from=0 /factor/core                         /usr/bin/core

ENTRYPOINT ["factor"]

CMD ["--version"]
