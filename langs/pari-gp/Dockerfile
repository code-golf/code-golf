FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl

ENV VER=2.17.0

RUN curl https://pari.math.u-bordeaux.fr/pub/pari/unix/pari-$VER.tar.gz | tar xz

WORKDIR /pari-$VER

RUN ./Configure        \
    --disable-mmap     \
    --graphic=none     \
    --mt=pthread       \
    --prefix=/usr      \
    --static           \
    --without-gmp      \
    --without-readline \
 && make -j`nproc` gp  \
 && make install       \
 && strip /usr/bin/gp

COPY pari-gp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pari-gp -s -static /pari-gp.c
 
FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gp              \
              /usr/bin/pari-gp         /usr/bin/

ENTRYPOINT ["gp"]

CMD ["--version-short"]
