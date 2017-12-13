FROM debian:stretch-slim

ENV CGO_ENABLED=0 GOROOT_BOOTSTRAP=/usr/local/go

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git make nasm

# https://golang.org/dl/
RUN curl -sSL https://storage.googleapis.com/golang/go1.10beta1.linux-amd64.tar.gz \
  | tar -xzC /usr/local

RUN git clone https://go.googlesource.com/go \
 && cd go                                    \
 && git checkout e7f95b3                     \
 && git config user.email a@b.c              \
 && git fetch https://go.googlesource.com/go refs/changes/97/82997/1 \
 && git cherry-pick FETCH_HEAD               \
 && cd src                                   \
 && ./make.bash                              \
 && chmod +rx /root

ENV GOCACHE=/tmp GOPATH=/root/go PATH=/go/bin:$PATH

RUN go get -d github.com/buildkite/terminal       \
 && cd /root/go/src/github.com/buildkite/terminal \
 && git checkout -q b0f19a1

RUN go get -d github.com/julienschmidt/httprouter       \
 && cd /root/go/src/github.com/julienschmidt/httprouter \
 && git checkout -q e1b9828

RUN go get -d github.com/lib/pq       \
 && cd /root/go/src/github.com/lib/pq \
 && git checkout -q 83612a5

RUN go get -d github.com/pmezard/go-difflib/difflib       \
 && cd /root/go/src/github.com/pmezard/go-difflib/difflib \
 && git checkout -q 792786c

CMD go build -ldflags -s -o app                    \
 && nasm -f bin -o run-container run-container.asm \
 && chmod +x run-container
