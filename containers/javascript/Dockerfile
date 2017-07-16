FROM alpine:edge

RUN apk add --no-cache curl g++ gcc make musl-dev perl python2

RUN curl https://ftp.mozilla.org/pub/firefox/releases/45.0.2esr/source/firefox-45.0.2esr.source.tar.xz \
  | tar xJf -

RUN cd /firefox-45.0.2esr/js/src                               \
 && SHELL=/bin/sh ./configure --disable-jemalloc --prefix=/usr \
 && make -j4                                                   \
 && make install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/js              /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   /usr/lib/
COPY --from=0 /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["/usr/bin/js"]
