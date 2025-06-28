FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl qbe scdoc

ENV VER=0.24.2

RUN curl -#L https://git.sr.ht/~sircmpwn/harec/archive/$VER.tar.gz \
  | tar xz --strip-components 1

RUN mv configs/linux.mk config.mk \
 && make install

RUN curl -# https://git.sr.ht/~sircmpwn/hare/archive/$VER.tar.gz | tar xz

WORKDIR /hare-$VER

RUN mv configs/linux.mk config.mk \
 && make && make install

COPY hare.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/hare -s /hare.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1   /lib/
COPY --from=0 /usr/bin/as                \
              /usr/bin/hare              \
              /usr/bin/ld                \
              /usr/bin/qbe               /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.44.so    \
              /usr/lib/libctf.so.0       \
              /usr/lib/libjansson.so.4   \
              /usr/lib/libsframe.so.1    \
              /usr/lib/libz.so.1         \
              /usr/lib/libzstd.so.1      /usr/lib/
COPY --from=0 /usr/local/bin/hare        \
              /usr/local/bin/harec       /usr/local/bin/
COPY --from=0 /usr/local/src/hare/stdlib /usr/local/src/hare/stdlib

ENTRYPOINT ["hare"]

CMD ["--version"]
