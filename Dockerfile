FROM debian:stretch

ENV CGO_ENABLED=0 GOPATH=/go PATH=/usr/local/go/bin:$PATH

WORKDIR /go

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl gcc git gnupg libc6-dev make openjdk-8-jre-headless vim-common

# https://golang.org/dl/
RUN curl -sSL https://storage.googleapis.com/golang/go1.9beta2.linux-amd64.tar.gz | tar -xzC /usr/local

RUN go get -d github.com/gorilla/handlers  \
 && cd /go/src/github.com/gorilla/handlers \
 && git checkout -q a4d79d4

RUN go get -d github.com/lib/pq  \
 && cd /go/src/github.com/lib/pq \
 && git checkout -q 8837942

RUN go get -d github.com/sergi/go-diff/... \
 && cd /go/src/github.com/sergi/go-diff    \
 && git checkout -q feef008

RUN go get -d github.com/tdewolff/minify  \
 && cd /go/src/github.com/tdewolff/minify \
 && git checkout -q 2d28d6e

RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -

RUN apt-get update && apt-get install -y --no-install-recommends nodejs

RUN npm install -g csso-cli@1.0.0 csso@3.1.1

RUN curl -L https://github.com/google/brotli/archive/v0.6.0.tar.gz \
  | tar xzf -                                                      \
 && cd brotli-0.6.0                                                \
 && make                                                           \
 && mv bin/bro /usr/local/bin

RUN curl http://dl.google.com/closure-compiler/compiler-20170521.tar.gz \
  | tar -zxf - -C /

# Bashisms FTW.
RUN ln -snf /bin/bash /bin/sh

CMD css=`cat static/*.css | csso /dev/stdin`                                            \
 &&  js=`java -jar /*.jar --assume_function_wrapper static/{codemirror{,-*},script}.js` \
 && echo -e "package main                                                             \n\
                                                                                      \n\
    const cssHash = \"`md5sum <<< "$css" | tr -d ' -'`\"                              \n\
    const  jsHash = \"`md5sum <<< "$js"  | tr -d ' -'`\"                              \n\
                                                                                      \n\
    var cssBr   = []byte{`bro     <<< "$css" | xxd -i`}                               \n\
    var cssGzip = []byte{`gzip -9 <<< "$css" | xxd -i`}                               \n\
    var  jsBr   = []byte{`bro     <<< "$js"  | xxd -i`}                               \n\
    var  jsGzip = []byte{`gzip -9 <<< "$js"  | xxd -i`}                               \n\
                                                                                      \n\
    " > static.go                   \
 && go build -ldflags '-s' -o app   \
 && gcc                             \
    -fno-asynchronous-unwind-tables \
    -nostdlib                       \
    -O2                             \
    -o run-container                \
    -s                              \
    -Wall                           \
    -Werror                         \
    -Wl,--build-id=none             \
    run-container.S run-container.c \
 && strip -R .comment run-container
