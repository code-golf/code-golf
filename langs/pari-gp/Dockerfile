FROM alpine:3.20 AS builder

# install base deps
RUN apk add --no-cache build-base curl gmp-dev readline-dev

# download pari
RUN curl https://pari.math.u-bordeaux.fr/pub/pari/unix/pari-2.15.5.tar.gz \
  | tar xz

# compile pari
RUN mv pari-* pari       \
 && cd pari              \
 && ./Configure          \
    --prefix=/usr        \
    --mt=pthread         \
 && make -j`nproc` gp    \
 && make install         \
 && make install-bin-sta \
 # dunno what is
    /usr/bin/gp

# configure gp (doesnt cooperate, copying for now)
#RUN touch ~/.gprc                           \
#    histfile="/root/.gp_history" >> ~/.gprc \
#    colors="darkbg"              >> ~/.gprc \
#    lines=40                     >> ~/.gprc \
#    parisizemax=4G               >> ~/.gprc \
#    thraedsizemax=1G             >> ~/.gprc \
#    read="/root/.gprc.gp"        >> ~/.gprc \
# && touch ~/.gprc.gp
COPY .gprc /

RUN touch ~/.gprc.gp

# download gp2c
RUN curl https://pari.math.u-bordeaux.fr/pub/pari/GP2C/gp2c-0.0.13.tar.gz \
  | tar xz

# compile gp2c
RUN mv gp2c-* gp2c                          \
 && cd gp2c                                 \
 && ./configure                             \
    --prefix=`pwd`/../usr                   \
    --with-paricfg=../usr/lib/pari/pari.cfg \
 && make check                              \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin                    /bin
COPY --from=0 /lib                    /lib
COPY --from=0 /usr/bin                /usr/bin
COPY --from=0 /usr/lib                /usr/lib
COPY --from=0 /usr/local/bin          /usr/local/bin
COPY --from=0 /usr/local/lib          /usr/local/lib

ENTRYPOINT ["gp2c-run"]

CMD ["--version"]
