FROM alpine:3.20 AS builder

ENV VER=0.58.5

RUN apk add --no-cache binutils-gold build-base clang curl

RUN curl -L https://github.com/ponylang/ponyup/archive/refs/tags/0.8.3.tar.gz \
  | tar xz --strip-components 1

RUN SHELL="sh" ./ponyup-init.sh \
 && /root/.local/share/ponyup/bin/ponyup update ponyc release-$VER

RUN mv /root/.local/share/ponyup/ponyc-* /usr/local/pony

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/rm /bin/sh                   /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/libz.so.1    /lib/
COPY --from=0 /usr/bin/clang /usr/bin/ld.gold            /usr/bin/
COPY --from=0 /usr/lib/Scrt1.o                           \
              /usr/lib/crti.o                            \
              /usr/lib/crtn.o                            \
              /usr/lib/libatomic.a                       \
              /usr/lib/libc.so                           \
              /usr/lib/libclang-cpp.so.17                \
              /usr/lib/libdl.a                           \
              /usr/lib/libffi.so.8                       \
              /usr/lib/libgcc_s.so                       \
              /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libjansson.so.4                   \
              /usr/lib/libLLVM-17.so                     \
              /usr/lib/liblzma.so.5                      \
              /usr/lib/libm.a                            \
              /usr/lib/libpthread.a                      \
              /usr/lib/librt.a                           \
              /usr/lib/libssp_nonshared.a                \
              /usr/lib/libstdc++.so.6                    \
              /usr/lib/libxml2.so.2                      \
              /usr/lib/libzstd.so.1                      \
              /usr/lib/gcc/*/*/crtbeginS.o               \
              /usr/lib/gcc/*/*/crtendS.o                 \
              /usr/lib/gcc/*/*/libgcc.a                  /usr/lib/
COPY --from=0 /usr/local/pony/bin/ponyc                  /usr/local/pony/bin/
COPY --from=0 /usr/local/pony/lib/x86-64/libponyrt-pic.a /usr/local/pony/lib/x86-64/
COPY --from=0 /usr/local/pony/packages                   /usr/local/pony/packages

COPY /pony /usr/bin/

ENTRYPOINT ["pony"]

CMD ["--version"]
