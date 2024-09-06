FROM alpine:3.20 AS builder

RUN apk add --no-cache curl

RUN curl -o fennel https://fennel-lang.org/downloads/fennel-1.5.1 \
 && chmod +x fennel

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/lua.
RUN sed -i '1c#!/usr/bin/lua' fennel

FROM codegolf/lang-lua

COPY --from=0 fennel /usr/bin/

ENTRYPOINT ["fennel"]

CMD ["-v"]
