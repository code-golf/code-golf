FROM ubuntu:18.04 as builder

RUN apt-get update && apt-get install -y curl libatomic1 libedit2 libc6-dev

RUN curl https://swift.org/builds/swift-5.1.5-release/ubuntu1804/swift-5.1.5-RELEASE/swift-5.1.5-RELEASE-ubuntu18.04.tar.gz \
  | tar --directory / --strip-components 1 -xzf -

RUN rm                     \
    /usr/bin/clangd        \
    /usr/bin/lldb-*        \
    /usr/bin/llvm-*        \
    /usr/bin/sourcekit-lsp \
    /usr/bin/swift-*

FROM scratch

COPY --from=0 / /

ENTRYPOINT [ "swift", "--version" ]
