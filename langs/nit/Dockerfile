FROM alpine:3.22 AS builder

RUN apk add --no-cache ccache curl dash gc-dev gcc libunwind-dev make musl-dev pkgconf

ENV VER=32b9c90

WORKDIR /usr/local

RUN curl -#L https://github.com/nitlang/nit/tarball/$VER \
  | tar xz --strip-components 1

RUN make -j`nproc` tools \
 && strip bin/nit

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                         /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1                   /lib/
COPY --from=0 /usr/bin/as /usr/bin/cc /usr/bin/ld   /usr/bin/
COPY --from=0 /usr/include                          /usr/include
COPY --from=0 /usr/lib/crti.o                       \
              /usr/lib/crtn.o                       \
              /usr/lib/libbfd-2.44.so               \
              /usr/lib/libc.so                      \
              /usr/lib/libctf.so.0                  \
              /usr/lib/libgc.so.1                   \
              /usr/lib/libgcc_s.so                  \
              /usr/lib/libgcc_s.so.1                \
              /usr/lib/libgmp.so.10                 \
              /usr/lib/libisl.so.23                 \
              /usr/lib/libjansson.so.4              \
              /usr/lib/liblzma.so.5                 \
              /usr/lib/libm.a                       \
              /usr/lib/libmpc.so.3                  \
              /usr/lib/libmpfr.so.6                 \
              /usr/lib/libsframe.so.1               \
              /usr/lib/libssp_nonshared.a           \
              /usr/lib/libunwind.so.8               \
              /usr/lib/libz.so.1                    \
              /usr/lib/libzstd.so.1                 /usr/lib/
COPY --from=0 /usr/lib/gcc/*/*/crtbeginS.o          \
              /usr/lib/gcc/*/*/crtendS.o            \
              /usr/lib/gcc/*/*/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/libexec/gcc/*/*/cc1              \
              /usr/libexec/gcc/*/*/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/local/bin/nit                    /usr/local/bin/
COPY --from=0 /usr/local/lib/core                   /usr/local/lib/core
COPY --from=0 /usr/local/src                        /usr/local/src
