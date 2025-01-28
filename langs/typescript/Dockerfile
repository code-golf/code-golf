FROM node:23.6.1-alpine3.21 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=5.7.3

RUN npm install --global typescript@$VER

COPY typescript.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/typescript -s -static typescript.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /usr/bin/env /usr/bin/typescript           /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libstdc++.so.6                    /usr/lib/
COPY --from=0 /usr/local/bin/node /usr/local/bin/tsc     /usr/local/bin/
COPY --from=0 /usr/local/lib/node_modules/typescript/lib /usr/local/lib

ENTRYPOINT ["typescript"]

CMD ["--version"]
