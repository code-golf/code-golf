FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes autoconf automake bash cmake g++ git libtool make texinfo

ENV MAKEFLAGS=-j8 VER=0.1.14

RUN git clone --progress --recurse-submodules -j4 -b v$VER \
    --shallow-submodules https://github.com/egel-lang/egel

WORKDIR /egel

RUN cp CMakeLists.static CMakeLists.txt

WORKDIR /egel/vendor

RUN bash makeall.sh

WORKDIR /egel/build

RUN cmake -DCMAKE_BUILD_TYPE='MinSizeRel' .. \
 && make install && strip /usr/local/bin/egel

COPY egel.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/egel -s /egel.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/libm.so.6      \
              /lib/x86_64-linux-gnu/libstdc++.so.6 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=0 /usr/bin/egel                        /usr/bin/
COPY --from=0 /usr/local/bin/egel                  /usr/local/bin/
COPY --from=0 /usr/local/lib/egel                  /usr/local/lib/egel

ENTRYPOINT ["egel"]

CMD ["--version"]
