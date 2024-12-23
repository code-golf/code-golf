FROM alpine:3.21 AS builder

ENV VER=0.58.7

RUN apk add --no-cache binutils-gold build-base clang curl dash

RUN curl -#L https://github.com/ponylang/ponyup/archive/refs/tags/0.8.4.tar.gz \
  | tar xz --strip-components 1

RUN SHELL='sh' ./ponyup-init.sh \
 && /root/.local/share/ponyup/bin/ponyup update ponyc release-$VER

RUN mv /root/.local/share/ponyup/ponyc-* /usr/local/pony

COPY pony.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pony -s -static pony.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                              /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /usr/bin/clang                             \
              /usr/bin/ld.gold                           \
              /usr/bin/pony                              /usr/bin/
COPY --from=0 /usr/lib/Scrt1.o                           \
              /usr/lib/crti.o                            \
              /usr/lib/crtn.o                            \
              /usr/lib/libatomic.a                       \
              /usr/lib/libc.so                           \
              /usr/lib/libclang-cpp.so.19.1              \
              /usr/lib/libdl.a                           \
              /usr/lib/libffi.so.8                       \
              /usr/lib/libgcc_s.so                       \
              /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libjansson.so.4                   \
              /usr/lib/libLLVM.so.19.1                   \
              /usr/lib/liblzma.so.5                      \
              /usr/lib/libm.a                            \
              /usr/lib/libpthread.a                      \
              /usr/lib/librt.a                           \
              /usr/lib/libssp_nonshared.a                \
              /usr/lib/libstdc++.so.6                    \
              /usr/lib/libxml2.so.2                      \
              /usr/lib/libz.so.1                         \
              /usr/lib/libzstd.so.1                      \
              /usr/lib/gcc/x86_64*/14.2.0/crtbeginS.o    \
              /usr/lib/gcc/x86_64*/14.2.0/crtendS.o      \
              /usr/lib/gcc/x86_64*/14.2.0/libgcc.a       /usr/lib/
COPY --from=0 /usr/local/pony/bin/ponyc                  /usr/local/pony/bin/
COPY --from=0 /usr/local/pony/lib/x86-64/libponyrt-pic.a /usr/local/pony/lib/x86-64/
COPY --from=0 /usr/local/pony/packages                   /usr/local/pony/packages

ENTRYPOINT ["pony"]

CMD ["--version"]
