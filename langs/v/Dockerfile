FROM alpine:3.22 AS builder

RUN apk add --no-cache dash gc-dev gcc git make musl-dev

ENV VER=0.4.12

WORKDIR /usr/local

RUN wget -O- https://github.com/vlang/v/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make \
 && strip v

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
