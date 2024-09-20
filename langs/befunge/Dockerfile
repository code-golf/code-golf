FROM alpine:3.20 AS builder

RUN apk add --no-cache curl



FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT [""]

CMD ["--version"]
