FROM alpine:3.21 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=0.1

RUN wget -cO /usr/bin/iogii https://golfscript.com/iogii/iogii-$VER-alpha \
 && chmod +x /usr/bin/iogii

COPY iogiiwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogiiwrapper -s -static iogiiwrapper.c

FROM codegolf/lang-ruby

COPY --from=0 /usr/bin/env          \
              /usr/bin/iogii        \
              /usr/bin/iogiiwrapper /usr/bin/

FROM codegolf/lang-base

COPY --from=1 . /

ENTRYPOINT ["iogiiwrapper"]

CMD ["--version"]
