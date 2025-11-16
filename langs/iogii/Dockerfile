FROM codegolf/lang-haskell

FROM alpine:3.22 AS builder

RUN apk add --no-cache gcc musl-dev

ENV BIN='/usr/local/bin/iogii' VER=1.1

RUN wget -qO $BIN https://golfscript.com/iogii/iogii-$VER \
 && chmod +x $BIN

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/ruby.
RUN sed -i '1c#!/usr/bin/ruby' $BIN

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s iogii.c

FROM codegolf/lang-ruby

COPY --from=0 /                    /
COPY --from=1 /usr/bin/iogii       /usr/bin/
COPY --from=1 /usr/local/bin/iogii /usr/local/bin/

ENTRYPOINT ["iogii"]

CMD ["--version"]
