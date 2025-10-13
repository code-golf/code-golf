FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc libncurses-dev libtinfo5 xz-utils

ENV VER=1.10.1

RUN curl -#L https://downloads.sourceforge.net/fbc/FreeBASIC-$VER-linux-x86_64.tar.xz \
  | tar xJ --strip-components 1

RUN ./install.sh -i /usr

COPY basic.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/basic -s basic.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                             /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libmvec.so.1    \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/libtinfo.so.5   \
              /lib/x86_64-linux-gnu/libz.so.1       /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/as                           \
              /usr/bin/basic                        \
              /usr/bin/fbc                          \
              /usr/bin/gcc                          \
              /usr/bin/ld                           /usr/bin/
COPY --from=0 /usr/include/freebasic                /usr/include/freebasic
COPY --from=0 /usr/lib/freebasic                    /usr/lib/freebasic
COPY --from=0 /usr/lib/gcc                          /usr/lib/gcc
COPY --from=0 /usr/lib/x86_64-linux-gnu             /usr/lib/x86_64-linux-gnu

ENTRYPOINT ["basic"]

CMD ["--version"]
