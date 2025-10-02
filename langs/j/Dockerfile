FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes clang git make nasm

ENV CC='clang' CFLAGS='-Wno-error' VER=9.6

RUN git clone --branch J$VER-beta3 --progress \
    https://github.com/jsoftware/jsource

WORKDIR /jsource/make2

RUN jplatform=linux j64x=j64avx2 ./build_all.sh

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/libm.so.6      /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=0 /jsource/jlibrary/addons/dev         /opt/j/addons/dev
COPY --from=0 /jsource/mpir/linux/x86_64/libgmp.so \
              /jsource/bin/linux/j64avx2/libj.so   \
              /jsource/jlibrary/bin/profile.ijs    /opt/j/bin/
COPY --from=0 /jsource/bin/linux/j64avx2/jconsole  /opt/j/bin/j
COPY --from=0 /jsource/jlibrary/system             /opt/j/system

ENTRYPOINT ["/opt/j/bin/j"]
