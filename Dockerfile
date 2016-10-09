FROM golang:1.7.1

ENV CGO_ENABLED 0

RUN go get -d github.com/tdewolff/minify  \
 && cd /go/src/github.com/tdewolff/minify \
 && git checkout -q a6728ce

RUN go get -d github.com/valyala/fasthttp  \
 && cd /go/src/github.com/valyala/fasthttp \
 && git checkout -q b43280d
