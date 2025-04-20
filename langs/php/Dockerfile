FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl

ENV CFLAGS='-O2 -flto' LDFLAGS='-O2 -flto -all-static' VER=8.4.6

RUN curl -#L https://php.net/distributions/php-$VER.tar.xz | tar xJ

WORKDIR /php-$VER

RUN LDFLAGS= ./configure       \
    --disable-all              \
    --disable-gcc-global-regs  \
    --disable-ipv6             \
    --disable-zend-signals     \
    --prefix=/usr              \
 && make -j`nproc` install-cli \
 && strip /usr/bin/php

RUN echo display_errors=stderr > /usr/lib/php.ini

FROM codegolf/lang-base

COPY --from=0 /usr/bin/php     /usr/bin/
COPY --from=0 /usr/lib/php.ini /usr/lib/

ENTRYPOINT ["php"]

CMD ["--version"]
