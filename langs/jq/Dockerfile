FROM alpine:3.22 AS builder

ENV BIN='/usr/bin/jq' VER=1.8.0

RUN wget -qO $BIN https://github.com/jqlang/jq/releases/download/jq-$VER/jq-linux64 \
 && chmod +x $BIN

FROM codegolf/lang-base

COPY --from=0 /usr/bin/jq /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
