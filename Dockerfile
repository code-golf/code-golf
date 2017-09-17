FROM debian:stretch

ENV CGO_ENABLED=0 GOPATH=/go PATH=/usr/local/go/bin:$PATH

WORKDIR /go

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl gcc git gnupg libc6-dev make nasm openjdk-8-jre-headless vim-common

# https://golang.org/dl/
RUN curl -sSL https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz | tar -xzC /usr/local

RUN go get -d github.com/buildkite/terminal  \
 && cd /go/src/github.com/buildkite/terminal \
 && git checkout -q c8b6c2b

RUN go get -d github.com/lib/pq  \
 && cd /go/src/github.com/lib/pq \
 && git checkout -q e422674

RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -

RUN apt-get update && apt-get install -y --no-install-recommends nodejs

RUN npm install -g csso-cli@1.0.0 csso@3.1.1

RUN curl -L https://github.com/google/brotli/archive/v0.6.0.tar.gz \
  | tar xzf -                                                      \
 && cd brotli-0.6.0                                                \
 && make                                                           \
 && mv bin/bro /usr/local/bin

RUN curl http://dl.google.com/closure-compiler/compiler-20170806.tar.gz \
  | tar -zxf - -C /

# Bashisms FTW.
RUN ln -snf /bin/bash /bin/sh

CMD css=`cat static/*.css | csso /dev/stdin`                                            \
 &&  js=`java -jar /*.jar --assume_function_wrapper static/{codemirror{,-*},script}.js` \
 && echo -e "package main                                                             \n\
                                                                                      \n\
    const cssPath = \"/`md5sum <<< "$css" | tr -d ' -'`\"                             \n\
    const  jsPath = \"/`md5sum <<< "$js"  | tr -d ' -'`\"                             \n\
                                                                                      \n\
    var cssBr   = []byte{`bro     <<< "$css" | xxd -i`}                               \n\
    var cssGzip = []byte{`gzip -9 <<< "$css" | xxd -i`}                               \n\
    var favicon = []byte{`xxd -i < static/favicon.ico`}                               \n\
    var  jsBr   = []byte{`bro     <<< "$js"  | xxd -i`}                               \n\
    var  jsGzip = []byte{`gzip -9 <<< "$js"  | xxd -i`}                               \n\
                                                                                      \n\
    " > static.go                                  \
 && go build -ldflags '-s' -o app                  \
 && nasm -f bin -o run-container run-container.asm \
 && chmod +x run-container
