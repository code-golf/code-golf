FROM alpine:3.22 AS builder

RUN apk add --no-cache clang curl make

# Also rebuild GolfScript, iogii.
ENV VER=3.4.7

WORKDIR /ruby

RUN curl -#L https://cache.ruby-lang.org/pub/ruby/3.4/ruby-$VER.tar.xz \
  | tar xJ --strip-components 1

# Prefer clang to gcc as it produces smaller binaries with adequate stack.
RUN ./configure            \
    --disable-install-doc  \
    --disable-jit-support  \
    --prefix=/usr          \
    --with-gcc=clang       \
 && make -j`nproc` install \
 && strip -s /usr/bin/ruby

# Remove some default gems we don't need.
RUN cd /usr/lib/ruby/3.4.0 && rm -rf \
                      English.rb     \
                      benchmark*     \
                      bundler*       \
    x86_64-linux-musl/bundler.so     \
                      coverage*      \
    x86_64-linux-musl/coverage.so    \
                      csv*           \
                      debug.rb       \
                      fiddle*        \
    x86_64-linux-musl/fiddle.so      \
                      mkmf.rb        \
                      objspace*      \
    x86_64-linux-musl/objspace.so    \
                      openssl*       \
    x86_64-linux-musl/openssl.so     \
                      psych*         \
    x86_64-linux-musl/psych.so       \
                      rdoc*          \
    x86_64-linux-musl/rdoc.so        \
                      ripper*        \
    x86_64-linux-musl/ripper.so

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/ruby       /usr/bin/
COPY --from=0 /usr/lib/ruby       /usr/lib/ruby

ENTRYPOINT ["ruby"]

CMD ["--version"]
