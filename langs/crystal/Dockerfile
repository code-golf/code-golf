FROM crystallang/crystal:1.13.2-alpine AS builder

ENV VER=1.13.2

RUN apk add --no-cache autoconf automake curl g++ libstdc++ libstdc++-dev libtool llvm15 llvm15-dev

ENV LLVM_CONFIG=/usr/lib/llvm15/bin/llvm-config

RUN curl https://www.hboehm.info/gc/gc_source/gc-8.2.4.tar.gz \
  | tar xz --strip-components 1

RUN ./autogen.sh          \
 && ./configure           \
    --disable-shared      \
    --enable-large-config \
    --enable-static       \
 && make && make install

RUN curl -L https://github.com/crystal-lang/crystal/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` crystal interpreter=1 release=1

# Added extra dep here for now because if the previous line initiates, it runs for over 1,100 seconds. No, thanks. Not now.
RUN apk add libevent-dev

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 .build/crystal           /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8     \
              /usr/lib/libgcc_s.so.1   \
              /usr/lib/libLLVM-15.so   \
              /usr/lib/liblzma.so.5    \
              /usr/lib/libpcre2-8.so.0 \
              /usr/lib/libstdc++.so.6  \
              /usr/lib/libxml2.so.2    \
              /usr/lib/libzstd.so.1    /usr/lib/

ENTRYPOINT ["crystal"]

CMD ["-v"]
