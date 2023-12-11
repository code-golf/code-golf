FROM alpine:3.19 AS builder

RUN apk add --no-cache curl

RUN curl -L https://github.com/babashka/babashka/releases/download/v1.3.186/babashka-1.3.186-linux-amd64-static.tar.gz | tar xz

FROM codegolf/lang-base

COPY --from=0 /bb /usr/bin/clojure

ENTRYPOINT ["clojure"]

CMD ["--version"]
