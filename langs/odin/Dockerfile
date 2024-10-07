FROM alpine:3.20 AS builder

RUN apk add --no-cache clang curl llvm-dev make

RUN curl -L https://github.com/odin-lang/Odin/archive/refs/tags/dev-2024-10.tar.gz | tar xz

RUN mv Odin-* /usr/local/tmp

RUN ln -s /usr/local/tmp/base /usr/local/bin/base \
 && ln -s /usr/local/tmp/core /usr/local/bin/core

RUN mkdir /usr/local/bin/vendor

WORKDIR /usr/local/tmp

RUN ODIN_ROOT=/usr/local/bin make release -j `nproc`

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/rm /bin/sh                                /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/libz.so.1                 /lib/
COPY --from=0 /usr/bin/clang /usr/bin/ld                              /usr/bin/
COPY --from=0 /usr/lib/crt1.o                                         \
              /usr/lib/crti.o                                         \
              /usr/lib/crtn.o                                         \
              /usr/lib/libbfd-2.42.so                                 \
              /usr/lib/libc.so                                        \
              /usr/lib/libclang-cpp.so.17                             \
              /usr/lib/libctf.so.0                                    \
              /usr/lib/libffi.so.8                                    \
              /usr/lib/libgcc_s.so                                    \
              /usr/lib/libgcc_s.so.1                                  \
              /usr/lib/libjansson.so.4                                \
              /usr/lib/libLLVM-17.so                                  \
              /usr/lib/liblzma.so.5                                   \
              /usr/lib/libm.a                                         \
              /usr/lib/libsframe.so.1                                 \
              /usr/lib/libssp_nonshared.a                             \
              /usr/lib/libstdc++.so.6                                 \
              /usr/lib/libxml2.so.2                                   \
              /usr/lib/libzstd.so.1                                   /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtbegin.o \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtend.o   \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/libgcc.a   /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/
COPY --from=0 /usr/local/tmp/odin                                     /usr/local/bin/
COPY --from=0 /usr/local/bin/base                                     /usr/local/bin/base
COPY --from=0 /usr/local/bin/core                                     /usr/local/bin/core
COPY --from=0 /usr/local/bin/vendor                                   /usr/local/bin/vendor

COPY /odinwrapper /usr/bin/

ENTRYPOINT ["odinwrapper"]

CMD ["--version"]
