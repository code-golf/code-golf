FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV LDFLAGS='-static' VER=2.17.3

RUN curl -#L https://pari.math.u-bordeaux.fr/pub/pari/unix/pari-$VER.tar.gz \
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

COPY pari-gp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pari-gp -s -static pari-gp.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/gp /usr/bin/pari-gp /usr/bin/

ENTRYPOINT ["pari-gp"]

CMD ["--version"]
