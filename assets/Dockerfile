FROM alpine:3.14

RUN apk add --no-cache brotli curl openjdk8-jre perl zopfli

# Closure Compiler
RUN curl http://dl.google.com/closure-compiler/compiler-20200614.tar.gz \
  | tar xz                                                              \
 && mv closure-compiler-v20200614.jar closure-compiler.jar              \
 && chmod +r closure-compiler.jar

WORKDIR /work

CMD ["assets/build"]
