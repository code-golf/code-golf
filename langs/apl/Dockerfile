FROM debian:trixie-slim AS builder

ENV VER=20.0.53963

ADD https://www.dyalog.com/uploads/php/download.dyalog.com/download.php?file=20.0/linux_64_${VER}_unicode.x86_64.deb /dyalog.deb

RUN dpkg -i dyalog.deb

FROM codegolf/lang-base

COPY --from=0 /bin/bash                                           /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                     \
              /lib/x86_64-linux-gnu/libdl.so.2                    \
              /lib/x86_64-linux-gnu/libm.so.6                     \
              /lib/x86_64-linux-gnu/libpthread.so.0               \
              /lib/x86_64-linux-gnu/librt.so.1                    \
              /lib/x86_64-linux-gnu/libtinfo.so.6                 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                         /lib64/
COPY --from=0 /opt/mdyalog/20.0/64/unicode/dyalog                 /opt/mdyalog/20.0/64/unicode/
COPY --from=0 /opt/mdyalog/20.0/64/unicode/aplkeys/default        \
              /opt/mdyalog/20.0/64/unicode/aplkeys/utf8           \
              /opt/mdyalog/20.0/64/unicode/aplkeys/utf8codes      /opt/mdyalog/20.0/64/unicode/aplkeys/
COPY --from=0 /opt/mdyalog/20.0/64/unicode/apltrans/default       \
              /opt/mdyalog/20.0/64/unicode/apltrans/utf8          /opt/mdyalog/20.0/64/unicode/apltrans/
COPY --from=0 /opt/mdyalog/20.0/64/unicode/scriptbin/dyalogscript /opt/mdyalog/20.0/64/unicode/scriptbin/
COPY --from=0 /usr/bin/dyalogscript                               /usr/bin/

ENTRYPOINT ["dyalogscript"]
