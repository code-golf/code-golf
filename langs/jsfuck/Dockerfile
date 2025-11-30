ARG VER=25.2.1

FROM node:$VER-alpine3.22 AS builder

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /usr/lib/libgcc_s.so.1  \
              /usr/lib/libstdc++.so.6 /usr/lib/
COPY --from=0 /usr/local/bin/node     /usr/local/bin/jsfuck

ENTRYPOINT ["jsfuck"]

CMD ["--version"]
