FROM alpine:3.22

ENV VER=0.36.0

RUN apk add --no-cache curl firefox ttf-dejavu \
 && curl -L https://github.com/mozilla/geckodriver/releases/download/v$VER/geckodriver-v$VER-linux64.tar.gz \
  | tar xz -C /usr/local/bin \
 && apk del --no-cache curl

CMD ["geckodriver", "--allow-hosts", "firefox", "--host", "0.0.0.0"]
