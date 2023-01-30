FROM debian:bullseye-slim as builder

RUN apt-get update \
 && apt-get install -y binutils curl libatomic1 libc6-dev libedit2 libz3-4

RUN curl https://download.swift.org/swift-5.7.3-release/ubuntu2004/swift-5.7.3-RELEASE/swift-5.7.3-RELEASE-ubuntu20.04.tar.gz \
  | tar xz --directory / --strip-components 1

RUN strip /usr/bin/swift

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu                    /lib/x86_64-linux-gnu
COPY --from=0 /lib64                                   /lib64
COPY --from=0 /usr/bin/swift                           /usr/bin/
COPY --from=0 /usr/include                             /usr/include
COPY --from=0 /usr/lib/clang                           /usr/lib/clang
COPY --from=0 /usr/lib/swift                           /usr/lib/swift
COPY --from=0 /usr/lib/x86_64-linux-gnu/libstdc++.so.6 \
              /usr/lib/x86_64-linux-gnu/libuuid.so.1   /usr/lib/x86_64-linux-gnu/

ENTRYPOINT ["swift"]

CMD ["--version"]
