FROM alpine:3.19 AS builder

RUN apk add --no-cache curl

ENV VER 1.3.190

RUN curl -L https://github.com/babashka/babashka/releases/download/v$VER/babashka-$VER-linux-amd64-static.tar.gz | tar xz

FROM codegolf/lang-base

COPY --from=0 /bb /usr/bin/clojure

ENTRYPOINT ["clojure"]

CMD ["--version"]
