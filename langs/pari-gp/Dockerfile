FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl https://pari.math.u-bordeaux.fr/pub/pari/unix/pari-2.17.0.tar.gz \
  | tar xz --strip-components 1

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

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sed /bin/sh /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/gp               /usr/bin/

COPY /pari-gp /usr/bin/

ENTRYPOINT ["pari-gp"]

CMD ["--version"]
