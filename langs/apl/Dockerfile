FROM debian:bookworm-slim AS builder

ARG RELEASE=18.2 \
    BUILD=45405 \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y curl dpkg libtinfo5

RUN curl -o dyalog.deb https://www.dyalog.com/uploads/php/download.dyalog.com/download.php?file=$RELEASE/linux_64_$RELEASE.${BUILD}_unicode.x86_64.deb \
 && dpkg -i dyalog.deb

FROM codegolf/lang-base

ARG DIR=18.2

COPY --from=0 /bin/bash /bin/cat /bin/chmod /bin/sh          /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                \
              /lib/x86_64-linux-gnu/libdl.so.2               \
              /lib/x86_64-linux-gnu/libm.so.6                \
              /lib/x86_64-linux-gnu/libpthread.so.0          \
              /lib/x86_64-linux-gnu/librt.so.1               \
              /lib/x86_64-linux-gnu/libtinfo.so.5            \
              /lib/x86_64-linux-gnu/libtinfo.so.6            /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                    /lib64/
COPY          /apl                                           /usr/bin/
COPY --from=0 /opt/mdyalog/$DIR/64/unicode/dyalog            /opt/mdyalog/$DIR/64/unicode/
COPY --from=0 /opt/mdyalog/$DIR/64/unicode/aplkeys/default   \
              /opt/mdyalog/$DIR/64/unicode/aplkeys/file      \
              /opt/mdyalog/$DIR/64/unicode/aplkeys/utf8codes /opt/mdyalog/$DIR/64/unicode/aplkeys/
COPY --from=0 /opt/mdyalog/$DIR/64/unicode/apltrans/default  \
              /opt/mdyalog/$DIR/64/unicode/apltrans/file     /opt/mdyalog/$DIR/64/unicode/apltrans/
COPY --from=0 /opt/mdyalog/$DIR/64/unicode/scriptbin         /opt/mdyalog/$DIR/64/unicode/scriptbin

RUN chmod +x /usr/bin/apl

ENTRYPOINT ["apl"]
