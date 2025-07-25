FROM alpine:3.20 AS builder

RUN apk add --no-cache curl

ENV VER=d3cc04f

RUN curl -#L https://github.com/m-ender/labyrinth/tarball/$VER \
  | tar xz --strip-components 1

FROM codegolf/lang-ruby

COPY --from=0 *.rb /

FROM codegolf/lang-base

COPY --from=1 . /

CMD ["ruby", "labyrinth.rb"]
