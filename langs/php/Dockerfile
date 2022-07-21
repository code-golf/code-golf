FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl

RUN curl -L https://php.net/distributions/php-7.4.30.tar.xz | tar xJ

ENV CFLAGS='-O2 -flto' LDFLAGS='-O2 -flto'

RUN cd php-7.4.30                                             \
 && ./configure                                               \
    --disable-all                                             \
    --disable-gcc-global-regs                                 \
    --disable-ipv6                                            \
    --disable-zend-signals                                    \
    --prefix=/usr                                             \
 && LDFLAGS="$LDFLAGS -all-static" make -j`nproc` install-cli \
 && strip /usr/bin/php

RUN echo display_errors=stderr > /usr/lib/php.ini

FROM codegolf/lang-base

COPY --from=0 /usr/bin/php     /usr/bin/
COPY --from=0 /usr/lib/php.ini /usr/lib/

ENTRYPOINT ["php"]

CMD ["-r", "echo phpversion();"]
