FROM crystallang/crystal:1.20.1-alpine AS builder

RUN apk add --no-cache dash \
 && apk del --no-cache libxml2-static python3

RUN strip /usr/bin/crystal

FROM codegolf/lang-base

# Can't use --parents because dash gets renamed, cc & pkg-config are symlinks.
COPY --from=0           /usr/bin/dash                        /bin/sh
COPY --from=0           /usr/bin/cc /usr/bin/pkg-config      /usr/bin/
COPY --from=0 --parents /lib/ld-musl-*.so.1                  \
                        /usr/bin/crystal                     \
                        /usr/include                         \
                        /usr/lib                             \
                        /usr/libexec/gcc                     \
                        /usr/share/crystal                   \
                        /usr/x86_64-alpine-linux-musl/bin/ld /

ENTRYPOINT ["crystal"]

CMD ["--version"]
