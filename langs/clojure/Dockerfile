FROM alpine:3.22 AS builder

ENV VER=1.12.208

RUN wget -O- https://github.com/babashka/babashka/releases/download/v$VER/babashka-$VER-linux-amd64-static.tar.gz \
  | tar xz --directory /usr/bin

FROM codegolf/lang-base

COPY --from=0 /usr/bin/bb /usr/bin/clojure

ENTRYPOINT ["clojure"]

CMD ["--version"]
