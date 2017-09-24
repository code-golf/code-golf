FROM debian:stretch

ENV CGO_ENABLED=0 GOPATH=/go PATH=/usr/local/go/bin:$PATH

WORKDIR /go

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl gcc git make nasm

# https://golang.org/dl/
RUN curl -sSL https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz \
  | tar -xzC /usr/local

RUN go get -d github.com/buildkite/terminal  \
 && cd /go/src/github.com/buildkite/terminal \
 && git checkout -q c8b6c2b

RUN go get -d github.com/lib/pq  \
 && cd /go/src/github.com/lib/pq \
 && git checkout -q e422674

CMD go build -ldflags '-s' -o app                  \
 && nasm -f bin -o run-container run-container.asm \
 && chmod +x run-container
