FROM alpine:3.22 AS builder

RUN apk add --no-cache bison cmake curl flex g++ gc-dev gc-static make

ENV LDFLAGS='-static' VER=a68b339

RUN curl -#L https://github.com/Wodan58/Moy/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /build

RUN cmake ..        \
 && cmake --build . \
 && strip joy

COPY joy.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/joywrapper -s -static /joy.c

FROM codegolf/lang-base

COPY --from=0 /build/joy /usr/bin/joywrapper /usr/bin/

ENTRYPOINT ["joywrapper"]

CMD ["--version"]
