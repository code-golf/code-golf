FROM node:23.1.0-alpine3.20 AS builder

RUN apk add --no-cache build-base

ENV VER=0.4.0

RUN npm install --prefix=/usr jsfuck@$VER

COPY jsfuck.c /

RUN gcc -s -o jsfuck jsfuck.c

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-x86_64.so.1           /lib/
COPY --from=0 /jsfuck /usr/bin/env               /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1             \
              /usr/lib/libstdc++.so.6            /usr/lib/
COPY --from=0 /usr/local/bin/node                /usr/local/bin/
COPY --from=0 /usr/node_modules/.bin             /usr/node_modules/bin
COPY --from=0 /usr/node_modules/jsfuck/fuck.js   \
              /usr/node_modules/jsfuck/jsfuck.js /usr/node_modules/jsfuck/

ENTRYPOINT ["jsfuck"]

CMD ["--version"]
