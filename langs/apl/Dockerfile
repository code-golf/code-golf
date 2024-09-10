FROM debian:bookworm-slim AS builder

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y curl dpkg libtinfo5

RUN curl -o dyalog.deb https://www.dyalog.com/uploads/php/download.dyalog.com/download.php?file=19.0/linux_64_19.0.50027_unicode.x86_64.deb \
 && dpkg -i dyalog.deb

FROM codegolf/lang-base

ENV DYALOG=/opt/mdyalog/19.0/64/unicode

COPY --from=0 /bin/bash /bin/cat /bin/sh            /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/librt.so.1      \
              /lib/x86_64-linux-gnu/libtinfo.so.5   \
              /lib/x86_64-linux-gnu/libtinfo.so.6   /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 $DYALOG/dyalog                        $DYALOG/
COPY --from=0 $DYALOG/aplkeys/default               \
              $DYALOG/aplkeys/utf8                  \
              $DYALOG/aplkeys/utf8codes             $DYALOG/aplkeys/
COPY --from=0 $DYALOG/apltrans/default              \
              $DYALOG/apltrans/utf8                 $DYALOG/apltrans/
COPY --from=0 $DYALOG/scriptbin/dyalogscript        /usr/bin/

COPY /apl /usr/bin/

ENTRYPOINT ["apl"]
