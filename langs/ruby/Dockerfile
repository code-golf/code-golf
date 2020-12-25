FROM alpine:3.12 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base clang curl openssl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/3.0/ruby-3.0.0.tar.xz \
  | tar xJf -

# CC=clang produces the same executable size as using CFLAGS=-flto with default gcc, but
# compiling with LTO with gcc leads to an issue where the stack size available for recursion
# is significantly reduced. It may be worth trying to enable LTO for Clang. That appears
# to require building the LLVMgold plugin.
RUN cd ruby-3.0.0          \
 && CC=clang ./configure   \
    --disable-install-doc  \
    --prefix=/usr          \
 && make -j`nproc` install \
 && strip -s /usr/bin/ruby

# Remove some default gems we don't need.
RUN cd /usr/lib/ruby/3.0.0 && rm -rf \
                      English.rb     \
                      benchmark      \
                      benchmark.rb   \
                      bundler        \
                      bundler.rb     \
    x86_64-linux-musl/bundler.so     \
                      coverage       \
                      coverage.rb    \
    x86_64-linux-musl/coverage.so    \
                      csv            \
                      csv.rb         \
                      debug.rb       \
                      fiddle         \
                      fiddle.rb      \
    x86_64-linux-musl/fiddle.so      \
                      mkmf.rb        \
                      objspace       \
                      objspace.rb    \
    x86_64-linux-musl/objspace.so    \
                      openssl        \
                      openssl.rb     \
    x86_64-linux-musl/openssl.so     \
                      psych          \
                      psych.rb       \
    x86_64-linux-musl/psych.so       \
                      rdoc           \
                      rdoc.rb        \
    x86_64-linux-musl/rdoc.so        \
                      ripper         \
                      ripper.rb      \
    x86_64-linux-musl/ripper.so

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin/ruby            /usr/bin/
COPY --from=0 /usr/lib/ruby/3.0.0      /usr/lib/ruby/3.0.0

ENTRYPOINT ["ruby", "-e", "puts RUBY_VERSION"]
