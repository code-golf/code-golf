FROM golang:1.12.4-alpine

ENV GOCACHE=/go/.go/cache GOPATH=/go/.go/path

RUN apk --no-cache add g++ git musl-dev

CMD ["go", "build", "-ldflags", "-extldflags -static -linkmode external -s"]
