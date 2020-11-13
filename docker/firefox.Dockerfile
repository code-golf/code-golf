FROM alpine:3.12

RUN apk add --no-cache curl firefox \
 && curl -L https://github.com/mozilla/geckodriver/releases/download/v0.28.0/geckodriver-v0.28.0-linux64.tar.gz \
  | tar xz -C /usr/local/bin        \
 && apk del --no-cache curl

CMD ["geckodriver", "--host", "0.0.0.0"]
