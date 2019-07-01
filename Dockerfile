FROM golang:1.13beta1-alpine

ENV GOBIN=/go GOCACHE=/go/.go/cache GOPATH=/go/.go/path TZ=Europe/London

RUN apk --no-cache add g++ git musl-dev \
 && GOBIN=/bin go get github.com/cespare/reflex

CMD ["go", "build", "-ldflags", "-extldflags -static -linkmode external -s"]
