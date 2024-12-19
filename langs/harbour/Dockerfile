FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc libpcre3-dev make zlib1g-dev

ENV VER=ce1bc20

RUN curl -#L https://github.com/harbour/core/tarball/$VER \
  | tar xz --strip-components 1

RUN HB_BUILD_CONTRIBS='no'   \
    HB_BUILD_STRIP='all'     \
    HB_INSTALL_PREFIX='/usr' \
    make -j`nproc` install

COPY harbourwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbourwrapper -s -static harbourwrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpcre.so.3    \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/libz.so.1       /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/harbour                      \
              /usr/bin/harbourwrapper               \
              /usr/bin/hbmk2                        /usr/bin/
COPY --from=0 /usr/include/harbour                  /usr/include/harbour

ENTRYPOINT ["harbourwrapper"]

CMD ["--version"]
