FROM alpine:3.21 AS builder

RUN apk add --no-cache dash gcc ghc musl-dev

ENV BIN='/usr/local/bin/iogii' VER=0.3

RUN wget -qO $BIN https://golfscript.com/iogii/iogii-$VER-beta \
 && chmod +x $BIN

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s iogii.c

FROM codegolf/lang-ruby

COPY --from=0 /usr/bin/dash        /bin/sh
COPY --from=0 /usr/bin/env         \
              /usr/bin/iogii       \
              /usr/bin/runghc      /usr/bin/
COPY --from=0 /usr/local/bin/iogii /usr/local/bin/

FROM codegolf/lang-haskell

# Compiling iogii will generate and execute Haskell code.
# See https://github.com/code-golf/code-golf/issues/1770.

FROM codegolf/lang-base

COPY --from=1 . /
COPY --from=2 . /

ENTRYPOINT ["iogii"]

CMD ["--version"]
