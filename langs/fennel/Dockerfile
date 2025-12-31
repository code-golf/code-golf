FROM alpine:3.23 AS builder

ENV VER=1.6.1

ADD --chmod=777 https://fennel-lang.org/downloads/fennel-$VER /usr/bin/fennel

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/lua.
RUN sed -i '1c#!/usr/bin/lua' /usr/bin/fennel

FROM codegolf/lang-lua

COPY --from=0 /usr/bin/fennel /usr/bin/

ENTRYPOINT ["fennel"]

CMD ["--version"]
