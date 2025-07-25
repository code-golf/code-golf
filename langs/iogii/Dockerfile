FROM alpine:3.22 AS builder

RUN apk add --no-cache dash gcc musl-dev

ENV BIN='/usr/local/bin/iogii' VER=1.1

RUN wget -qO $BIN https://golfscript.com/iogii/iogii-$VER \
 && chmod +x $BIN

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/ruby.
RUN sed -i '1c#!/usr/bin/ruby' $BIN

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s iogii.c

FROM codegolf/lang-ruby

COPY --from=0 /usr/bin/dash        /bin/sh
COPY --from=0 /usr/bin/env         \
              /usr/bin/iogii       /usr/bin/
COPY --from=0 /usr/local/bin/iogii /usr/local/bin/

FROM codegolf/lang-haskell

# Compiling iogii will generate and execute Haskell code.
# See https://github.com/code-golf/code-golf/issues/1770.

FROM codegolf/lang-base

COPY --from=1 . /
COPY --from=2 . /

ENTRYPOINT ["iogii"]

CMD ["--version"]
