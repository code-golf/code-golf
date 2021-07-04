FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl

RUN curl -L https://php.net/distributions/php-7.4.20.tar.xz | tar xJ

ENV CFLAGS='-O2 -flto -ffat-lto-objects' \
   LDFLAGS='-O2 -flto -ffat-lto-objects'

RUN cd php-7.4.20                                             \
 && ./configure                                               \
    --disable-all                                             \
    --disable-gcc-global-regs                                 \
    --disable-ipv6                                            \
    --disable-zend-signals                                    \
    --prefix=/usr                                             \
 && LDFLAGS="$LDFLAGS -all-static" make -j`nproc` install-cli \
 && strip /usr/bin/php

RUN echo display_errors=stderr > /usr/lib/php.ini

FROM scratch

COPY --from=0 /empty           /proc
COPY --from=0 /empty           /tmp
COPY --from=0 /usr/bin/php     /usr/bin/
COPY --from=0 /usr/lib/php.ini /usr/lib/

ENTRYPOINT ["php"]

CMD ["-r", "echo phpversion();"]
