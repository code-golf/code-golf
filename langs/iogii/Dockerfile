FROM codegolf/lang-haskell

FROM alpine:3.23 AS builder

RUN apk add --no-cache gcc musl-dev

ENV BIN='/usr/local/bin/iogii' VER=1.1.1

ADD --chmod=755 https://golfscript.com/iogii/iogii-$VER $BIN

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/ruby.
RUN sed -i '1c#!/usr/bin/ruby' $BIN

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s iogii.c

FROM codegolf/lang-ruby

ENV GHC='x86_64-linux-ghc-9.14.1-d531'

COPY --from=1 /usr/bin/gcc /usr/bin/iogii       /usr/bin/
COPY --from=0 /usr/lib                          /usr/lib
COPY --from=0 /usr/local/bin                    /usr/local/bin
COPY --from=1 /usr/local/bin/iogii              /usr/local/bin/
COPY --from=0 /usr/local/lib/settings           /usr/local/lib/
COPY --from=0 /usr/local/lib/package.conf.d     /usr/local/lib/package.conf.d
COPY --from=0 /usr/local/lib/$GHC/libHS*-*-*.*  /usr/local/lib/$GHC/
COPY --from=0 /usr/local/lib/$GHC/base*         /usr/local/lib/$GHC/base-4.22.0.0-d466
COPY --from=0 /usr/local/lib/$GHC/containers*   /usr/local/lib/$GHC/containers-0.8-3e81
COPY --from=0 /usr/local/lib/$GHC/ghc-internal* /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29
COPY --from=0 /usr/local/lib/$GHC/rts*          /usr/local/lib/$GHC/rts-1.0.3

ENTRYPOINT ["iogii"]

CMD ["--version"]
