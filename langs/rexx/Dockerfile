FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes gcc sysv-rc-conf wget

ENV VER=3.9.7

RUN wget -O libregina.deb https://sourceforge.net/projects/regina-rexx/files/regina-rexx/$VER/libregina3_$VER-1-amd64-Debian-12.deb \
 && dpkg -i libregina.deb

RUN wget -O regina.deb https://sourceforge.net/projects/regina-rexx/files/regina-rexx/$VER/regina-rexx_$VER-1-amd64-Debian-12.deb \
 && dpkg -i regina.deb

COPY rexx.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rexx -s rexx.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6          \
              /lib/x86_64-linux-gnu/libcrypt.so.1      \
              /lib/x86_64-linux-gnu/libdl.so.2         \
              /lib/x86_64-linux-gnu/libpthread.so.0    /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2              /lib64/
COPY --from=0 /usr/bin/regina /usr/bin/rexx            /usr/bin/
COPY --from=0 /usr/lib/x86_64-linux-gnu/libregina.so.3 /usr/lib/

ENTRYPOINT ["rexx"]

CMD ["--version"]
