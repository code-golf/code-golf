FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc make

ENV VER=3.8.5

RUN curl -#L https://jmvdveer.home.xs4all.nl/algol68g-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure   \
    --prefix=/usr \
 && make install  \
 && strip /usr/bin/a68g

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6        \
              /lib/x86_64-linux-gnu/libm.so.6        \
              /lib/x86_64-linux-gnu/libquadmath.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2            /lib64/
COPY --from=0 /usr/bin/a68g                          /usr/bin/algol-68

ENTRYPOINT ["algol-68"]

CMD ["--version"]
