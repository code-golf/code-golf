FROM debian:stretch-slim

ENV CGO_ENABLED=0 GOROOT_BOOTSTRAP=/usr/local/go

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git make nasm

# https://golang.org/dl/
RUN curl -sSL https://storage.googleapis.com/golang/go1.10beta2.linux-amd64.tar.gz \
  | tar -xzC /usr/local

RUN git clone https://go.googlesource.com/go \
 && cd go                                    \
 && git checkout 9f31353                     \
 && git config user.email a@b.c              \
 && git fetch https://go.googlesource.com/go refs/changes/97/82997/1 \
 && git cherry-pick FETCH_HEAD               \
 && cd src                                   \
 && ./make.bash                              \
 && chmod +rx /root

ENV GOCACHE=/tmp GOPATH=/root/go PATH=/go/bin:$PATH

CMD go build -ldflags -s -o app                    \
 && nasm -f bin -o run-container run-container.asm \
 && chmod +x run-container
