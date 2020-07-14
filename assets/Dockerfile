FROM alpine:3.12

RUN apk add --no-cache bash curl g++ make musl-dev nodejs-npm openjdk8-jre perl

# Brotli
RUN curl -L https://github.com/google/brotli/archive/v1.0.7.tar.gz \
  | tar xzf -                                                      \
 && cd brotli-1.0.7                                                \
 && make -j`nproc`                                                 \
 && mv bin/brotli /usr/bin

# Zopfli
RUN curl -L https://github.com/google/zopfli/archive/zopfli-1.0.3.tar.gz \
  | tar xzf -                                                            \
 && cd zopfli-zopfli-1.0.3                                               \
 && make -j`nproc`                                                       \
 && mv zopfli /usr/bin

# Closure Compiler
RUN curl http://dl.google.com/closure-compiler/compiler-20200112.tar.gz \
  | tar -zxf - -C /                                                     \
 && mv closure-compiler-v20200112.jar closure-compiler.jar              \
 && chmod +r closure-compiler.jar

# CSSO
RUN npm install -g csso-cli@3.0.0 csso@4.0.2

# Bashisms FTW.
RUN ln -snf /bin/bash /bin/sh

WORKDIR /work

CMD ["assets/build"]
