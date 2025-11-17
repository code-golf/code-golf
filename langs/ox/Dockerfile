FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes gcc musl-dev

ENV VER=9.30

ADD https://www.doornik.com/download/oxmetrics9/Ox_Console/oxconsole_$VER-0_amd64.deb ox.deb

RUN dpkg -i ox.deb

COPY ox.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/ox -s ox.c

FROM codegolf/lang-base

COPY --chown=nobody                                      \
     --from=0 /tmp                                       /dev/shm
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6            \
              /lib/x86_64-linux-gnu/libdl.so.2           \
              /lib/x86_64-linux-gnu/libpthread.so.0      \
              /lib/x86_64-linux-gnu/librt.so.1           /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                /lib64/
COPY --from=0 /usr/bin/ox                                /usr/bin/
COPY --from=0 /usr/share/OxMetrics9/ox/oxl               /usr/local/bin/
COPY --from=0 /usr/share/OxMetrics9/ox/include/oxstd.oxh /usr/local/include/
COPY --from=0 /usr/share/OxMetrics9/ox/libomp.so.5       \
              /usr/share/OxMetrics9/ox/libox.so.9        /usr/local/lib/

ENTRYPOINT ["ox"]
