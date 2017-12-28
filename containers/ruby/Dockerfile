FROM alpine:edge

RUN apk add --no-cache curl gcc linux-headers make musl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/2.5/ruby-2.5.0.tar.xz \
  | tar xJf -

RUN cd ruby-2.5.0         \
 && ./configure           \
    --disable-install-doc \
    --prefix=/usr         \
 && make                  \
 && make install

RUN find /usr/lib/ruby -name '*.md' -delete

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/ruby            /usr/bin/
COPY --from=0 /usr/lib/ruby            /usr/lib/ruby

ENTRYPOINT ["/usr/bin/ruby", "-e", "puts RUBY_VERSION"]
