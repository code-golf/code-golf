FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl

RUN curl -o fennel https://fennel-lang.org/downloads/fennel-1.5.1 \
 && chmod +x fennel

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/lua.
RUN sed -i '1c#!/usr/bin/lua' fennel

COPY fennelwrapper.c /

RUN gcc -static -s -o fennelwrapper fennelwrapper.c

FROM codegolf/lang-lua

COPY --from=0 fennel fennelwrapper /usr/bin/

ENTRYPOINT ["fennelwrapper"]

CMD ["-v"]
