FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl g++ make

ENV VER=0.26.1

RUN curl -#L https://github.com/IppClub/YueScript/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install \
 && strip /usr/local/bin/yue

COPY yuescript.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/yuescript -s yuescript.c

FROM codegolf/lang-lua

COPY --from=0 /usr/bin/yuescript \
              /usr/local/bin/yue /usr/bin/

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/libm.so.6      \
              /lib/x86_64-linux-gnu/libstdc++.so.6 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=1 /usr/bin                             /usr/bin

ENTRYPOINT ["yuescript"]

CMD ["--version"]
