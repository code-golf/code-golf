FROM golang:1.8.1

ENV CGO_ENABLED 0

RUN go get -d github.com/gorilla/handlers  \
 && cd /go/src/github.com/gorilla/handlers \
 && git checkout -q 13d7309

RUN go get -d github.com/tdewolff/minify  \
 && cd /go/src/github.com/tdewolff/minify \
 && git checkout -q 18372f3

RUN apt-get update \
 && apt-get install -y --no-install-recommends nodejs-legacy npm vim-common

RUN npm install -g csso-cli@1.0.0 csso@3.1.1

# Bashisms FTW.
RUN ln -snf /bin/bash /bin/sh

CMD css=`cat static/*.css | csso /dev/stdin`                \
 &&  js=`cat static/{codemirror{,-perl},script}.js`         \
 && echo -e "package main                                 \n\
    const cssHash = \"`md5sum <<< "$css" | tr -d ' -'`\"  \n\
    const  jsHash = \"`md5sum <<< "$js"  | tr -d ' -'`\"  \n\
                                                          \n\
    var cssGzip = []byte{`echo "$css" | gzip -9 | xxd -i`}\n\
    var  jsGzip = []byte{`echo "$js"  | gzip -9 | xxd -i`}\n\
    " > static.go                                           \
 && go build -ldflags '-s' -o app
