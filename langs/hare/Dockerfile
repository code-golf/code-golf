FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl qbe scdoc

ENV VERSION=0.24.2

RUN curl https://git.sr.ht/~sircmpwn/harec/archive/$VERSION.tar.gz | tar xz

RUN cd harec-$VERSION             \
 && mv configs/linux.mk config.mk \
 && make install

RUN curl https://git.sr.ht/~sircmpwn/hare/archive/$VERSION.tar.gz | tar xz

RUN cd hare-$VERSION              \
 && mv configs/linux.mk config.mk \
 && make && make install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/rm /bin/sh   /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1   \
              /lib/libz.so.1             /lib/
COPY --from=0 /usr/bin/as                \
              /usr/bin/ld                \
              /usr/bin/qbe               /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.42.so    \
              /usr/lib/libctf.so.0       \
              /usr/lib/libjansson.so.4   \
              /usr/lib/libsframe.so.1    \
              /usr/lib/libzstd.so.1      /usr/lib/
COPY --from=0 /usr/local/bin/hare        \
              /usr/local/bin/harec       /usr/local/bin/
COPY --from=0 /usr/local/src/hare/stdlib /usr/local/src/hare/stdlib

COPY /harewrapper /usr/bin/

ENTRYPOINT ["harewrapper"]
