FROM node:16.11.0-alpine3.14 as builder

ENV VERSION=2.1.6

RUN mkdir /empty

RUN apk add --no-cache binutils

RUN npm install -g @defasm/cli@1.2.2 @defasm/core@$VERSION

# Strip node.
RUN strip /usr/local/bin/node

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/node.
RUN sed -i '1c#!/usr/bin/node' /usr/local/lib/node_modules/@defasm/cli/cli.js

# Put the symlink in a folder to avoid docker COPY resolving it.
RUN mkdir /symlink \
 && ln -s /usr/lib/node_modules/@defasm/cli/cli.js /symlink/defasm

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1            \
              /lib/libz.so.1                      /lib/
COPY --from=0 /empty                              /proc
COPY --from=0 /empty                              /tmp
COPY --from=0 /symlink                            /usr/bin
COPY --from=0 /usr/local/bin/node                 /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1              \
              /usr/lib/libstdc++.so.6             /usr/lib/
COPY --from=0 /usr/local/lib/node_modules/@defasm /usr/lib/node_modules/@defasm

ENTRYPOINT ["/usr/bin/defasm"]
