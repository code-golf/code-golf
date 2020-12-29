FROM alpine:3.12 as builder

ENV VERSION=f749e59

RUN mkdir /empty

RUN apk add --no-cache build-base clang curl openssl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/3.0/ruby-3.0.0.tar.xz \
  | tar xJf -

RUN cd ruby-3.0.0          \
 && CC=clang ./configure   \
    --disable-install-doc  \
    --disable-jit-support  \
    --prefix=/usr          \
 && make -j`nproc` install \
 && strip -s /usr/bin/ruby
 
 
RUN curl -L https://github.com/m-ender/hexagony/tarball/$VERSION \
  | tar xz
  

COPY interpreter.rb   /m-ender-hexagony-$VERSION/

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1   /lib/
COPY --from=0 /empty                     /proc
COPY --from=0 /empty                     /tmp
COPY --from=0 /usr/bin/ruby              /usr/bin/
COPY --from=0 /usr/lib/ruby/3.0.0        /usr/lib/ruby/3.0.0
COPY --from=0 /m-ender-hexagony-f749e59  /usr/bin/hexagony

