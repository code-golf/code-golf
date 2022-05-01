FROM alpine:3.15 as builder

RUN apk add --no-cache build-base clang curl openssl-dev

RUN curl https://cache.ruby-lang.org/pub/ruby/3.1/ruby-3.1.2.tar.xz | tar xJ

# Clang produces the same size binary as using LTO with GCC, but GCC+LTO leads
# to an issue where the stack size available for recursion is significantly
# reduced. It may be worth trying to enable LTO for Clang. That appears to
# require building the LLVMgold plugin.
RUN cd ruby-3.1.2          \
 && CC=clang ./configure   \
    --disable-install-doc  \
    --disable-jit-support  \
    --prefix=/usr          \
 && make -j`nproc` install \
 && strip -s /usr/bin/ruby

# Remove some default gems we don't need.
RUN cd /usr/lib/ruby/3.1.0 && rm -rf \
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

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/ruby            /usr/bin/
COPY --from=0 /usr/lib/ruby            /usr/lib/ruby

ENTRYPOINT ["ruby"]

CMD ["-e", "puts RUBY_VERSION"]
