FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc libpcre3-dev make zlib1g-dev

ENV HB_BUILD_CONTRIBS='no' HB_BUILD_STRIP='all' VER=3.0.0

RUN curl -#L https://sourceforge.net/projects/harbour-project/files/source/$VER/harbour-$VER.tar.gz/download \
  | tar xz --strip-components 1

RUN make -j`nproc` install

WORKDIR /harbour

RUN curl -#L https://github.com/harbour/core/archive/master.tar.gz \
  | tar xz --strip-components 1

RUN HB_INSTALL_PREFIX='/usr' make -j`nproc` install

COPY harbour.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbour -s /harbour.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpcre.so.3    \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/libz.so.1       /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/harbour                      /usr/bin/
COPY --from=0 /usr/bin/hbmk2                        \
              /usr/local/bin/harbour                /usr/local/bin/
COPY --from=0 /usr/local/include/harbour            /usr/local/include/harbour

ENTRYPOINT ["harbour"]

CMD ["--version"]
