FROM alpine:3.17

RUN apk add --no-cache curl firefox ttf-dejavu \
 && curl -L https://github.com/mozilla/geckodriver/releases/download/v0.32.2/geckodriver-v0.32.2-linux64.tar.gz \
  | tar xz -C /usr/local/bin        \
 && apk del --no-cache curl

CMD ["geckodriver", "--allow-hosts", "firefox", "--host", "0.0.0.0"]
