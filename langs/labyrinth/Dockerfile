FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=d3cc04f

RUN curl -#L https://github.com/m-ender/labyrinth/tarball/$VER \
  | tar xz --directory /usr/lib --strip-components 1

COPY labyrinth.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/labyrinth -s labyrinth.c

FROM codegolf/lang-ruby

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /usr/bin/labyrinth      /usr/bin/
COPY --from=0 /usr/lib/direction.rb   \
              /usr/lib/interpreter.rb \
              /usr/lib/labyrinth.rb   \
              /usr/lib/point2d.rb     /usr/lib/

ENTRYPOINT ["labyrinth"]
