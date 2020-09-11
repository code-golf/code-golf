FROM alpine:3.12

RUN apk add --no-cache bash curl g++ make musl-dev openjdk8-jre perl

# Brotli
RUN curl -L https://github.com/google/brotli/archive/v1.0.9.tar.gz \
  | tar xzf -                                                      \
 && cd brotli-1.0.9                                                \
 && make -j`nproc`                                                 \
 && mv bin/brotli /usr/bin && cd / && rm -r brotli-1.0.9

# Zopfli
RUN curl -L https://github.com/google/zopfli/archive/zopfli-1.0.3.tar.gz \
  | tar xzf -                                                            \
 && cd zopfli-zopfli-1.0.3                                               \
 && make -j`nproc`                                                       \
 && mv zopfli /usr/bin && cd / && rm -r zopfli-zopfli-1.0.3

# Closure Compiler
RUN curl http://dl.google.com/closure-compiler/compiler-20200614.tar.gz \
  | tar -zxf -                                                          \
 && mv closure-compiler-v20200614.jar closure-compiler.jar              \
 && chmod +r closure-compiler.jar

WORKDIR /work

CMD ["assets/build"]
