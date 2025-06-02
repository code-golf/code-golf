FROM alpine:3.22 AS builder

ENV VER=1.12.200

RUN wget -O- https://github.com/babashka/babashka/releases/download/v$VER/babashka-$VER-linux-amd64-static.tar.gz | tar xz

FROM codegolf/lang-base

COPY --from=0 /bb /usr/bin/clojure

ENTRYPOINT ["clojure"]

CMD ["--version"]
