FROM crystallang/crystal:1.18.2-alpine AS builder

RUN apk add --no-cache dash \
 && apk del --no-cache libxml2-static python3

RUN strip /usr/bin/crystal

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                        /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1                  \
              /lib/libz.so.1                       /lib/
COPY --from=0 /usr/bin/cc                          \
              /usr/bin/crystal                     \
              /usr/bin/pkg-config                  /usr/bin/
COPY --from=0 /usr/include                         /usr/include
COPY --from=0 /usr/lib                             /usr/lib
COPY --from=0 /usr/libexec/gcc                     /usr/libexec/gcc
COPY --from=0 /usr/share/crystal                   /usr/share/crystal
COPY --from=0 /usr/x86_64-alpine-linux-musl/bin/ld /usr/x86_64-alpine-linux-musl/bin/

ENTRYPOINT ["crystal"]

CMD ["--version"]
