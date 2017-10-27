FROM alpine:edge

RUN apk add --no-cache bash curl gcc make musl-dev nodejs-npm openjdk8-jre perl

# Brotli
RUN curl -L https://github.com/google/brotli/archive/v1.0.1.tar.gz \
  | tar xzf -                                                      \
 && cd brotli-1.0.1                                                \
 && make -j`nproc`                                                 \
 && mv bin/brotli /usr/bin

# Closure Compiler
RUN curl http://dl.google.com/closure-compiler/compiler-20171023.tar.gz \
  | tar -zxf - -C /                                                     \
 && mv closure-compiler-v20171023.jar closure-compiler.jar

# CSSO & SVGO
RUN npm install -g csso-cli@1.1.0 csso@3.3.1 svgo@0.7.2

# Bashisms FTW.
RUN ln -snf /bin/bash /bin/sh

WORKDIR /work

CMD ["assets/build"]
