FROM alpine:3.20 AS builder

RUN apk add --no-cache clang curl llvm-dev make

RUN curl -L https://github.com/odin-lang/Odin/archive/refs/tags/dev-2024-08.tar.gz | tar xz

RUN mv Odin-* /usr/local/dist

WORKDIR /usr/local/dist

RUN make -j`nproc` release

# The final image gets close to exceeding 400 MiB prior to this hacky approach.
# It'll trick the compiler into thinking the `enormous` libraries are present.
# As a matter of fact, they would never be used for the purposes on code.golf.
RUN rm -rf vendor && mkdir vendor

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/chmod /bin/rm /bin/sh                     /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/libz.so.1                 /lib/
COPY           odinwrapper                                            /usr/bin/
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
COPY --from=0 /usr/local/dist/odin                                    /usr/local/dist/
COPY --from=0 /usr/local/dist/base                                    /usr/local/dist/base
COPY --from=0 /usr/local/dist/core                                    /usr/local/dist/core
COPY --from=0 /usr/local/dist/vendor                                  /usr/local/dist/vendor

RUN chmod +x /usr/bin/odinwrapper

ENTRYPOINT ["odinwrapper"]

CMD ["version"]
