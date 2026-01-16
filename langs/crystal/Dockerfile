ARG VER=1.18.2

FROM crystallang/crystal:$VER-alpine AS builder

ARG VER

RUN apk add --no-cache curl dash g++ llvm-dev

RUN curl -#L https://github.com/crystal-lang/crystal/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make interpreter=1 release=1 -j`nproc` install \
 && strip /usr/local/bin/crystal

COPY crystal.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/crystal -s crystal.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash            /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1      \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/bin/cc              \
              /usr/bin/crystal         \
              /usr/bin/pkg-config      /usr/bin/
COPY --from=0 /usr/lib/libLLVM-17.so   \
              /usr/lib/libffi.so.8     \
              /usr/lib/libgc.so.1      \
              /usr/lib/libgcc_s.so.1   \
              /usr/lib/liblzma.so.5    \
              /usr/lib/libpcre2-8.so.0 \
              /usr/lib/libpkgconf.so.5 \
              /usr/lib/libstdc++.so.6  \
              /usr/lib/libxml2.so.2    \
              /usr/lib/libzstd.so.1    /usr/lib/
COPY --from=0 /usr/local/bin/crystal   /usr/local/bin/
COPY --from=0 /usr/local/share/crystal /usr/local/share/crystal

ENTRYPOINT ["crystal"]

CMD ["--version"]
