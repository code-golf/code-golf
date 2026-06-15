FROM alpine:3.24

RUN apk add --no-cache firefox ttf-dejavu

ENV VER=0.37.0

ADD --unpack https://github.com/mozilla/geckodriver/releases/download/v$VER/geckodriver-v$VER-linux64.tar.gz /usr/local/bin/

CMD ["geckodriver", "--allow-hosts", "firefox", "--host", "0.0.0.0"]
