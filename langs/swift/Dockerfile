FROM ubuntu:20.04 as builder

RUN mkdir /empty

RUN apt-get update \
 && apt-get install -y binutils curl libatomic1 libc6-dev libedit2 libz3-4

RUN curl https://swift.org/builds/swift-5.3.3-release/ubuntu2004/swift-5.3.3-RELEASE/swift-5.3.3-RELEASE-ubuntu20.04.tar.gz \
  | tar xz --directory / --strip-components 1

RUN rm                              \
    /usr/bin/clangd                 \
    /usr/bin/lldb-*                 \
    /usr/bin/llvm-*                 \
    /usr/bin/sourcekit-lsp          \
    /usr/bin/swift-*                \
    /usr/lib/liblldb.so.10.0.0git   \
    /usr/lib/libsourcekitdInProc.so \
 && strip /usr/bin/clang-10 /usr/bin/swift

FROM scratch

COPY --from=0 /      /
COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp

ENTRYPOINT ["swift", "--version"]
