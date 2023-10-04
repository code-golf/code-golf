FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl

ENV VER 8.2.11

RUN curl -L https://php.net/distributions/php-$VER.tar.xz | tar xJ

ENV CFLAGS='-O2 -flto' LDFLAGS='-O2 -flto -all-static'

RUN cd php-$VER                \
 && LDFLAGS= ./configure       \
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

CMD ["-r", "echo phpversion();"]
