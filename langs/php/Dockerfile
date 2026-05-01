FROM alpine:3.23 AS builder

RUN apk add --no-cache autoconf bison gcc make musl-dev re2c

ENV CFLAGS='-O2 -flto' LDFLAGS='-O2 -flto -all-static' VER=8.5.5

WORKDIR /php

ADD https://github.com/php/php-src.git#php-$VER .

RUN ./buildconf --force        \
 && LDFLAGS= ./configure       \
    --disable-all              \
    --disable-gcc-global-regs  \
    --disable-ipv6             \
    --disable-zend-signals     \
    --prefix=/usr              \
 && make -j`nproc` install-cli \
 && strip /usr/bin/php

FROM codegolf/lang-base

COPY --from=0 /usr/bin/php /usr/bin/

COPY <<EOF /usr/lib/php.ini
display_errors=stderr
EOF

ENTRYPOINT ["php"]

CMD ["--version"]
