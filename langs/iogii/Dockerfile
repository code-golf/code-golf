FROM alpine:3.21 AS builder

RUN apk add --no-cache gcc musl-dev

ENV BIN='/usr/local/bin/iogii' VER=0.2

RUN wget -cO $BIN https://golfscript.com/iogii/iogii-$VER-alpha \
 && chmod +x $BIN

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s -static iogii.c

FROM codegolf/lang-ruby

COPY --from=0 /usr/bin/env   \
              /usr/bin/iogii /usr/bin/
COPY --from=0 /usr/local/bin /usr/local/bin

FROM codegolf/lang-base

COPY --from=1 . /

ENTRYPOINT ["iogii"]

CMD ["--version"]
