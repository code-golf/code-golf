FROM alpine:3.22 AS builder

ENV BIN='/usr/bin/fennel' VER=1.5.3

RUN wget -qO $BIN https://fennel-lang.org/downloads/fennel-$VER \
 && chmod +x $BIN

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/lua.
RUN sed -i '1c#!/usr/bin/lua' $BIN

FROM codegolf/lang-lua

COPY --from=0 /usr/bin/fennel /usr/bin/

ENTRYPOINT ["fennel"]

CMD ["--version"]
