FROM alpine:3.10 as builder

RUN apk add --no-cache curl gcc linux-headers make musl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/2.6/ruby-2.6.4.tar.xz \
  | tar xJf -

RUN cd ruby-2.6.4         \
 && ./configure           \
    --disable-install-doc \
    --prefix=/usr         \
 && make install          \
 && strip -s /usr/bin/ruby

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/ruby            /usr/bin/
COPY --from=0 /usr/lib/ruby/2.6.0      /usr/lib/ruby/2.6.0

ENTRYPOINT ["/usr/bin/ruby", "-e", "puts RUBY_VERSION"]
