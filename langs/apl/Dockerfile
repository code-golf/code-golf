FROM debian:bookworm-slim AS builder

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y curl dpkg libtinfo5

RUN curl -o dyalog.deb https://www.dyalog.com/uploads/php/download.dyalog.com/download.php?file=18.2/linux_64_18.2.45405_unicode.x86_64.deb \
 && dpkg -i dyalog.deb

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                \
              /lib/x86_64-linux-gnu/libdl.so.2               \
              /lib/x86_64-linux-gnu/libm.so.6                \
              /lib/x86_64-linux-gnu/libpthread.so.0          \
              /lib/x86_64-linux-gnu/librt.so.1               \
              /lib/x86_64-linux-gnu/libtinfo.so.5            /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                    /lib64/
COPY --from=0 /opt/mdyalog/18.2/64/unicode/dyalog            /usr/bin/
COPY --from=0 /opt/mdyalog/18.2/64/unicode/aplkeys/default   \
              /opt/mdyalog/18.2/64/unicode/aplkeys/file      \
              /opt/mdyalog/18.2/64/unicode/aplkeys/utf8codes /usr/bin/aplkeys/
COPY --from=0 /opt/mdyalog/18.2/64/unicode/apltrans/default  \
              /opt/mdyalog/18.2/64/unicode/apltrans/file     /usr/bin/apltrans/

ENTRYPOINT ["dyalog"]
