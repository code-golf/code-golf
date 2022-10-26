FROM alpine:3.16 as builder

RUN wget -O - 'https://github.com/elm/compiler/releases/download/0.19.1/binary-for-linux-64-bit.gz' \
    | gunzip -c >/usr/local/bin/elm
RUN chmod +x /usr/local/bin/elm
RUN apk add --no-cache nodejs=16.17.1-r0


RUN mkdir -p /elm/src

COPY elm.json /elm
COPY main.js /elm
COPY Pre.elm /elm/src
COPY Main.elm /elm/src
RUN  cd /elm && /usr/local/bin/elm make src/Pre.elm --optimize --output pre.js

COPY elm.sh /usr/bin/

RUN chmod -R 777 /elm


ENTRYPOINT ["elm.sh"]

CMD ["--version"]
