FROM node:21-alpine3.18 AS builder

RUN apk add --no-cache binutils
RUN strip /usr/local/bin/node

COPY index.js version.js package.json package-lock.json /typescript/
WORKDIR /typescript
RUN npm i

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1            \
              /lib/libz.so.1                      /lib/
COPY --from=0 /usr/local/bin/node                 /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1              \
              /usr/lib/libstdc++.so.6             /usr/lib/
COPY --from=0 /typescript                         /typescript

WORKDIR /typescript
ENTRYPOINT ["node"]
CMD ["version.js"]