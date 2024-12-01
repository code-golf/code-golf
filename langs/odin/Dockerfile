FROM alpine:3.20 AS builder

RUN apk add --no-cache clang curl llvm-dev make

ENV VER=2024-11

RUN curl -#L https://github.com/odin-lang/Odin/archive/refs/tags/dev-$VER.tar.gz \
  | tar xz --directory /usr/local/lib --strip-components 1

RUN ln -s /usr/local/lib/base \
          /usr/local/bin/base \
 && ln -s /usr/local/lib/core \
          /usr/local/bin/core \
 && mkdir /usr/local/bin/vendor

WORKDIR /usr/local/lib

RUN make -j`nproc` release \
 && strip /usr/local/lib/odin

COPY odinwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/odinwrapper -s -static /odinwrapper.c

FROM codegolf/lang-base

COPY --from=0 /bin/sh                     /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    \
              /lib/libz.so.1              /lib/
COPY --from=0 /usr/bin/clang              \
              /usr/bin/ld                 \
              /usr/bin/odinwrapper        /usr/bin/
COPY --from=0 /usr/lib/crt1.o             \
              /usr/lib/crti.o             \
              /usr/lib/crtn.o             \
              /usr/lib/libbfd-2.42.so     \
              /usr/lib/libc.so            \
              /usr/lib/libclang-cpp.so.17 \
              /usr/lib/libctf.so.0        \
              /usr/lib/libffi.so.8        \
              /usr/lib/libgcc_s.so        \
              /usr/lib/libgcc_s.so.1      \
              /usr/lib/libjansson.so.4    \
              /usr/lib/libLLVM-17.so      \
              /usr/lib/liblzma.so.5       \
              /usr/lib/libm.a             \
              /usr/lib/libsframe.so.1     \
              /usr/lib/libssp_nonshared.a \
              /usr/lib/libstdc++.so.6     \
              /usr/lib/libxml2.so.2       \
              /usr/lib/libzstd.so.1       /usr/lib/
COPY --from=0 /usr/local/lib/odin         /usr/local/bin/
COPY --from=0 /usr/local/bin/base         /usr/local/bin/base
COPY --from=0 /usr/local/bin/core         /usr/local/bin/core
COPY --from=0 /usr/local/bin/vendor       /usr/local/bin/vendor

COPY --from=0 /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtbegin.o \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtend.o   \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/libgcc.a   \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/

ENTRYPOINT ["odinwrapper"]

CMD ["--version"]
