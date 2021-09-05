FROM node:16.8.0-alpine3.14 as builder

ENV VERSION=2.1.3

RUN mkdir /empty

RUN apk add --no-cache g++ linux-headers make python3

RUN npm install -g caxa @defasm/cli@1.0.5 @defasm/core@$VERSION

RUN caxa -i /usr/local/lib/node_modules/@defasm/cli -o defasm -- \
    "{{caxa}}/node_modules/.bin/node" "{{caxa}}/main.js"

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1  \
              /lib/libz.so.1            /lib/
COPY --from=0 /empty                    /proc
COPY --from=0 /empty                    /tmp
COPY --from=0 /defasm /usr/bin/ld       /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.35.2.so \
              /usr/lib/libctf.so.0      \
              /usr/lib/libgcc_s.so.1    \
              /usr/lib/libstdc++.so.6   /usr/lib/

ENTRYPOINT ["/usr/bin/defasm"]
