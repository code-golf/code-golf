FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc musl-dev

# FIXME Bump to latest version.
ENV VER=8.02

RUN curl -#o ox.deb https://www.doornik.com/download/oxmetrics8/Ox_Console/oxcons_$VER-0_amd64.deb \
 && dpkg -i ox.deb

COPY ox.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/ox -s ox.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                 /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6           \
              /lib/x86_64-linux-gnu/libdl.so.2          \
              /lib/x86_64-linux-gnu/libm.so.6           /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2               /lib64/
COPY --from=0 /usr/bin/ox                               /usr/bin/
COPY --from=0 /usr/share/OxMetrics8/ox/bin64/libox.so.8 \
              /usr/share/OxMetrics8/ox/bin64/oxl.bin    /usr/local/bin/
COPY --from=0 /usr/share/OxMetrics8/ox/include          /usr/local/include

ENTRYPOINT ["ox"]
