FROM alpine:3.22 AS builder

RUN apk add --no-cache g++

ENV VER=1.0.1

RUN wget -O- https://www.golfscript.com/iogii/iogii-$VER.tgz \
  | tar xz --strip-components 2

COPY iogii.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s -static iogii.c

FROM codegolf/lang-ruby

COPY --from=0 /*.rb          /
COPY --from=0 /usr/bin/iogii /usr/bin/

ENTRYPOINT ["iogii"]

CMD ["--version"]
