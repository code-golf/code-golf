FROM alpine:edge as builder

RUN apk add --no-cache curl gcc make musl-dev

RUN curl -L https://php.net/distributions/php-7.3.6.tar.xz | tar xJf -

ENV CFLAGS='-O2 -flto' \
   LDFLAGS='-O2 -flto'

RUN cd php-7.3.6                                              \
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

COPY --from=0 /usr/bin/php     /usr/bin/
COPY --from=0 /usr/lib/php.ini /usr/lib/

ENTRYPOINT ["/usr/bin/php", "-r", "echo phpversion();"]
