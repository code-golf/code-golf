FROM debian:bookworm-slim AS builder

ENV VER=0.1.14

RUN apt-get update                   \
 && DEBIAN_FRONTEND="noninteractive" \
    apt-get install --yes clang cmake curl libfmt-dev make pkg-config

RUN curl https://ftp.gnu.org/gnu/lightning/lightning-2.2.3.tar.gz | tar xz

RUN cd lightning-2.2.3 \
 && ./configure        \
    --prefix="/usr"    \
 && make               \
 && make install

RUN curl -L https://github.com/egel-lang/egel/archive/refs/tags/v$VER.tar.gz | tar xz

RUN cd egel-$VER \
 && mkdir build  \
 && cd build     \
 && cmake ..     \
 && make         \
 && make install

FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT ["egel"]

CMD ["--version"]
