FROM alpine:3.21 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV BIN='/usr/local/bin/fennel' VER=1.5.3

RUN curl -#o $BIN https://fennel-lang.org/downloads/fennel-$VER \
 && chmod +x $BIN

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/lua.
RUN sed -i '1c#!/usr/bin/lua' $BIN

COPY fennel.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/fennel -s fennel.c

FROM codegolf/lang-lua

COPY --from=0 /usr/bin/fennel       /usr/bin/
COPY --from=0 /usr/local/bin/fennel /usr/local/bin/

ENTRYPOINT ["fennel"]

CMD ["--version"]
