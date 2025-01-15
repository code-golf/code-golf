FROM alpine:3.21 AS builder

ENV VER=d3cc04f

RUN wget -O- https://github.com/m-ender/labyrinth/tarball/$VER \
  | tar xz --strip-components 1

FROM codegolf/lang-ruby

COPY --from=0 *.rb /

FROM codegolf/lang-base

COPY --from=1 / /

CMD ["ruby", "labyrinth.rb"]
