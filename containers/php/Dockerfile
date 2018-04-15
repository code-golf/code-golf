FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl -L https://php.net/distributions/php-7.2.4.tar.xz | tar xJf -

RUN cd php-7.2.4   \
 && ./configure    \
    --disable-all  \
    --disable-ipv6 \
    --enable-cli   \
    --prefix=/usr  \
 && make -j`nproc` install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/php             /usr/bin/

ENTRYPOINT ["/usr/bin/php", "-r", "echo phpversion();"]
