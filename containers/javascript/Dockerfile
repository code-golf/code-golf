FROM alpine:edge

ENV SHELL /bin/sh

RUN apk add --no-cache \
    autoconf2.13 curl g++ gcc linux-headers make musl-dev python2

RUN curl https://ftp.mozilla.org/pub/firefox/releases/58.0/source/firefox-58.0.source.tar.xz \
  | tar xJf -

RUN cd /firefox-58.0/js/src                      \
 && ./configure --disable-jemalloc --prefix=/usr \
 && make -j4                                     \
 && make install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/js58            /usr/bin/js
COPY --from=0 /usr/lib/libgcc_s.so.1   /usr/lib/
COPY --from=0 /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["/usr/bin/js", "-v"]
