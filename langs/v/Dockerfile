FROM debian:buster-slim as builder

RUN mkdir /empty

RUN apt-get update && apt-get install -y gcc git make

RUN git clone -b weekly.2021.26 git://github.com/vlang/v /opt/v \
 && cd /opt/v && make && strip v

FROM debian:buster-slim

RUN apt-get update && apt-get install -y libc-dev

COPY --from=0 /opt/v/thirdparty/tcc /opt/v/thirdparty/tcc
COPY --from=0 /opt/v/v              /opt/v/v
COPY --from=0 /opt/v/vlib           /opt/v/vlib

FROM scratch

COPY --from=1 /bin   /bin
COPY --from=1 /lib   /lib
COPY --from=1 /lib64 /lib64
COPY --from=1 /opt/v /opt/v
COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY --from=1 /usr   /usr

COPY v /usr/bin/

ENTRYPOINT ["/usr/bin/v"]

CMD ["version"]
