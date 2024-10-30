ARG VER=1.14.0

FROM crystallang/crystal:$VER-alpine AS builder

ARG VER

RUN apk add --no-cache curl libstdc++-dev llvm18-dev llvm18-static zstd-dev zstd-static

RUN curl -L https://github.com/crystal-lang/crystal/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` crystal interpreter=1 release=1 static=1 \
 && strip /.build/crystal

FROM codegolf/lang-base

COPY --from=0 /bin/sh                               /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1              \
              /lib/libz.so.1                        /lib/
COPY --from=0 /.build/crystal                       \
              /usr/bin/cc                           \
              /usr/bin/ld                           \
              /usr/bin/pkg-config                   /usr/bin/
COPY --from=0 /usr/lib/Scrt1.o                      \
              /usr/lib/crti.o                       \
              /usr/lib/crtn.o                       \
              /usr/lib/libbfd-2.42.so               \
              /usr/lib/libc.so                      \
              /usr/lib/libctf.so.0                  \
              /usr/lib/libdl.a                      \
              /usr/lib/libevent.a                   \
              /usr/lib/libgc.a                      \
              /usr/lib/libgcc_s.so                  \
              /usr/lib/libgcc_s.so.1                \
              /usr/lib/libjansson.so.4              \
              /usr/lib/libpkgconf.so.5              \
              /usr/lib/libpthread.a                 \
              /usr/lib/librt.a                      \
              /usr/lib/libsframe.so.1               \
              /usr/lib/libssp_nonshared.a           \
              /usr/lib/libzstd.so.1                 \
              /usr/lib/gcc/*/*/crtbeginS.o          \
              /usr/lib/gcc/*/*/crtendS.o            \
              /usr/lib/gcc/*/*/libgcc.a             /usr/lib/
COPY --from=0 /usr/libexec/gcc/*/*/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/
COPY --from=0 /usr/share/crystal/src                /usr/share/crystal/src

ENTRYPOINT ["crystal"]

CMD ["--version"]
