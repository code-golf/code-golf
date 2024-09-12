FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl gmp-dev readline-dev

RUN curl https://pari.math.u-bordeaux.fr/pub/pari/unix/pari-2.15.5.tar.gz \
  | tar xz

RUN mv pari-* pari    \
 && cd pari           \
 && ./Configure       \
    --graphic=none    \
    --mt=pthread      \
    --prefix=/usr     \
    --static          \
 && make -j`nproc` gp \
 && make install      \
 && cp misc/gprc.dft $HOME/.gprc

RUN curl https://pari.math.u-bordeaux.fr/pub/pari/GP2C/gp2c-0.0.13.tar.gz \
  | tar xz

RUN mv gp2c-* gp2c                          \
 && cd gp2c                                 \
 && ./configure                             \
    --prefix=`pwd`/../usr                   \
    --with-paricfg=../usr/lib/pari/pari.cfg \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/chmod /bin/sh /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /usr/bin/gp                 /usr/bin/
COPY --from=0 /usr/lib/libgmp.so.10       \
              /usr/lib/libncursesw.so.6   \
              /usr/lib/libreadline.so.8   /usr/lib/

COPY /pari-gp /usr/bin/

# Remove
RUN chmod +x /usr/bin/pari-gp

ENTRYPOINT ["pari-gp"]

CMD ["--version"]
