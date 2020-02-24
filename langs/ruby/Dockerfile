FROM alpine:3.11 as builder

RUN apk add --no-cache build-base curl openssl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/2.7/ruby-2.7.0.tar.xz \
  | tar xJf -

RUN cd ruby-2.7.0         \
 && ./configure           \
    --disable-install-doc \
    --prefix=/usr         \
 && make install          \
 && strip -s /usr/bin/ruby

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/ruby            /usr/bin/
COPY --from=0 /usr/lib/ruby/2.7.0      /usr/lib/ruby/2.7.0

ENTRYPOINT ["/usr/bin/ruby", "-e", "puts RUBY_VERSION"]
