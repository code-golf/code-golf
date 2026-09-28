FROM alpine:3.23 AS builder

RUN apk add --no-cache dash gc-dev gcc git make musl-dev

# Pin the vc commit to the commit message that matches the V version.
#   e.g. [v:master] xxx - V $VER
ENV VER=0.5.1 VC_VER=f461dfe

WORKDIR /usr/local

ADD https://github.com/vlang/v.git#$VER .

# TODO Why no worky? - ADD https://github.com/vlang/vc.git?commit=$VC_VER vc
RUN git clone --depth=1000 https://github.com/vlang/vc \
 && cd vc && git reset --hard $VC_VER

RUN make GITCLEANPULL=true && strip v

COPY v.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/V -s -static /v.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash             /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/V                /usr/bin/
COPY --from=0 /usr/include              /usr/include
COPY --from=0 /usr/lib/crt1.o           \
              /usr/lib/crti.o           \
              /usr/lib/crtn.o           \
              /usr/lib/libatomic.so.1   \
              /usr/lib/libc.so          \
              /usr/lib/libdl.a          \
              /usr/lib/libm.a           \
              /usr/lib/libpthread.a     /usr/lib/
COPY --from=0 /usr/local/v              /usr/local/
COPY --from=0 /usr/local/thirdparty/tcc /usr/local/thirdparty/tcc
COPY --from=0 /usr/local/vlib           /usr/local/vlib

ENTRYPOINT ["V"]

CMD ["--version"]
