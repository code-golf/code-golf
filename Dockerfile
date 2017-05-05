FROM golang:1.8.1

ENV CGO_ENABLED 0

RUN go get -d github.com/gorilla/handlers  \
 && cd /go/src/github.com/gorilla/handlers \
 && git checkout -q 13d7309

RUN go get -d github.com/tdewolff/minify  \
 && cd /go/src/github.com/tdewolff/minify \
 && git checkout -q 18372f3
