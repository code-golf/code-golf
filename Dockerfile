FROM debian:stretch-slim

ENV CGO_ENABLED=0 GOROOT_BOOTSTRAP=/usr/local/go

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git make nasm

# https://golang.org/dl/
RUN curl -SSL https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz \
  | tar -xzC /usr/local

RUN git clone https://go.googlesource.com/go \
 && cd go                                    \
 && git checkout 9428023                     \
 && cd src                                   \
 && ./make.bash                              \
 && chmod +rx /root

ENV GOCACHE=/tmp GOPATH=/root/go PATH=/go/bin:$PATH

CMD go build -ldflags -s -o app                    \
 && nasm -f bin -o run-container run-container.asm \
 && chmod +x run-container
