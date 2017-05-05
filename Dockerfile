FROM golang:1.8.1

ENV CGO_ENABLED 0

RUN go get -d github.com/tdewolff/minify  \
 && cd /go/src/github.com/tdewolff/minify \
 && git checkout -q 18372f3
