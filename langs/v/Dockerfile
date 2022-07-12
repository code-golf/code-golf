FROM debian:bullseye-slim as builder

RUN apt-get update && apt-get install -y gcc git make

RUN git clone -b weekly.2022.22 https://github.com/vlang/v /opt/v \
 && cd /opt/v && make && strip v

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y libc-dev

COPY --from=0 /opt/v/thirdparty/tcc /opt/v/thirdparty/tcc
COPY --from=0 /opt/v/v              /opt/v/v
COPY --from=0 /opt/v/vlib           /opt/v/vlib

FROM codegolf/lang-base

COPY --from=1 /bin/cat /bin/rm /bin/sh                   /bin/
COPY --from=1 /lib/x86_64-linux-gnu                      /lib/x86_64-linux-gnu
COPY --from=1 /lib64                                     /lib64
COPY --from=1 /opt/v                                     /opt/v
COPY --from=1 /usr/include                               /usr/include
COPY --from=0 /usr/lib/x86_64-linux-gnu/crt1.o           \
              /usr/lib/x86_64-linux-gnu/crti.o           \
              /usr/lib/x86_64-linux-gnu/crtn.o           \
              /usr/lib/x86_64-linux-gnu/libatomic.so.1   \
              /usr/lib/x86_64-linux-gnu/libc.so          \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a \
              /usr/lib/x86_64-linux-gnu/libm.so          /usr/lib/x86_64-linux-gnu/

COPY v /usr/bin/

ENTRYPOINT ["v"]

CMD ["version"]
