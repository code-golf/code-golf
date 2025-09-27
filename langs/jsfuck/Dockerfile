ARG VER=24.4.1

FROM node:$VER-alpine3.22 AS builder

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/bin/node      /usr/bin/jsfuck
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["jsfuck"]

CMD ["--version"]
