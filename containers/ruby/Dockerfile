FROM alpine:edge

RUN apk add --no-cache curl gcc linux-headers make musl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/2.4/ruby-2.4.1.tar.xz \
  | tar xJf -

RUN cd ruby-2.4.1         \
 && ./configure           \
    --disable-install-doc \
    --prefix=/usr         \
 && make                  \
 && make install
