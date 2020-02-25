FROM golang:1.14.0-alpine

ENV CGO_ENABLED=0 GOBIN=/go GOCACHE=/go/.go/cache GOPATH=/go/.go/path TZ=Europe/London

RUN apk --no-cache add git && GOBIN=/bin go get github.com/cespare/reflex

CMD ["go", "build", "-ldflags", "-s"]
