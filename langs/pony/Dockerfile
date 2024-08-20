FROM alpine:3.20 AS builder

RUN apk add --no-cache binutils-gold build-base clang curl

RUN curl -L https://github.com/ponylang/ponyup/archive/refs/tags/0.8.3.tar.gz | tar xz --strip-components 1

RUN SHELL=sh ./ponyup-init.sh \
 && /root/.local/share/ponyup/bin/ponyup update ponyc release-0.58.5

RUN mv /root/.local/share/ponyup/ponyc-* /usr/local/pony

FROM codegolf/lang-base

COPY --from=0 / /

COPY --from=0 /lib/ld-musl-x86_64.so.1    \
              /lib/libz.so.1              /lib/
COPY --from=0 /usr/bin/clang              /usr/bin/
COPY --from=0 /usr/lib/libclang-cpp.so.17 \
              /usr/lib/libffi.so.8        \
              /usr/lib/libgcc_s.so.1      \
              /usr/lib/libLLVM-17.so      \
              /usr/lib/liblzma.so.5       \
              /usr/lib/libstdc++.so.6     \
              /usr/lib/libxml2.so.2       \
              /usr/lib/libzstd.so.1       /usr/lib/
COPY --from=0 /usr/local/pony             /usr/local/pony

COPY /pony /usr/bin/

RUN chmod +x /usr/bin/pony

ENTRYPOINT ["pony"]

CMD ["--version"]
