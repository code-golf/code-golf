FROM alpine:3.16

RUN apk add --no-cache curl firefox ttf-dejavu \
 && curl -L https://github.com/mozilla/geckodriver/releases/download/v0.29.1/geckodriver-v0.29.1-linux64.tar.gz \
  | tar xz -C /usr/local/bin        \
 && apk del --no-cache curl

CMD ["geckodriver", "--host", "0.0.0.0"]
