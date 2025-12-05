FROM alpine:3.22 AS builder

RUN apk add --no-cache binutils-gold clang curl dash

ENV SHELL='dash' VER=0.60.4

RUN curl -#L https://github.com/ponylang/ponyup/archive/master.tar.gz \
  | tar xz --strip-components 1

RUN ./ponyup-init.sh

RUN ~/.local/share/ponyup/bin/ponyup update ponyc release-$VER

RUN mv ~/.local/share/ponyup/ponyc-* /usr/local/pony

COPY pony.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pony -s pony.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                              /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1                        /lib/
COPY --from=0 /usr/bin/clang                             \
              /usr/bin/ld.gold                           \
              /usr/bin/pony                              /usr/bin/
COPY --from=0 /usr/lib/Scrt1.o                           \
              /usr/lib/crti.o                            \
              /usr/lib/crtn.o                            \
              /usr/lib/libLLVM.so.*                      \
              /usr/lib/libatomic.a                       \
              /usr/lib/libc.so                           \
              /usr/lib/libclang-cpp.so.*                 \
              /usr/lib/libdl.a                           \
              /usr/lib/libffi.so.8                       \
              /usr/lib/libgcc_s.so                       \
              /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libjansson.so.4                   \
              /usr/lib/liblzma.so.5                      \
              /usr/lib/libm.a                            \
              /usr/lib/libpthread.a                      \
              /usr/lib/librt.a                           \
              /usr/lib/libssp_nonshared.a                \
              /usr/lib/libstdc++.so.6                    \
              /usr/lib/libxml2.so.2                      \
              /usr/lib/libz.so.1                         \
              /usr/lib/libzstd.so.1                      \
              /usr/lib/gcc/*/14.2.0/crtbeginS.o          \
              /usr/lib/gcc/*/14.2.0/crtendS.o            \
              /usr/lib/gcc/*/14.2.0/libgcc.a             /usr/lib/
COPY --from=0 /usr/local/pony/bin/ponyc                  /usr/local/bin/
COPY --from=0 /usr/local/pony/lib/x86-64/libponyrt-pic.a /usr/local/lib/
COPY --from=0 /usr/local/pony/packages                   /usr/local/packages

ENTRYPOINT ["pony"]

CMD ["--version"]
