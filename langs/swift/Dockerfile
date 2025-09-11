FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes binutils curl libatomic1 libc6-dev libedit2 libncurses6 libsqlite3-0 libz3-4

ENV VER=6.1.3

RUN curl -#L https://download.swift.org/swift-$VER-release/debian12/swift-$VER-RELEASE/swift-$VER-RELEASE-debian12.tar.gz \
  | tar xz --strip-components 1

RUN strip /usr/bin/swift /usr/bin/swift-frontend \
 && rm -r /usr/lib/swift/embedded                \
 && find /usr/lib/swift -name '*.so' -exec strip {} +

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6           \
              /lib/x86_64-linux-gnu/libdl.so.2          \
              /lib/x86_64-linux-gnu/libgcc_s.so.1       \
              /lib/x86_64-linux-gnu/libm.so.6           \
              /lib/x86_64-linux-gnu/libncurses.so.6     \
              /lib/x86_64-linux-gnu/libpthread.so.0     \
              /lib/x86_64-linux-gnu/librt.so.1          \
              /lib/x86_64-linux-gnu/libtinfo.so.6       \
              /lib/x86_64-linux-gnu/libutil.so.1        \
              /lib/x86_64-linux-gnu/libz.so.1           /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2               /lib64/
COPY --from=0 /usr/bin/swift /usr/bin/swift-frontend    /usr/bin/
COPY --from=0 /usr/include                              /usr/include
COPY --from=0 /usr/lib/x86_64-linux-gnu/libsqlite3.so.0 \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6  \
              /usr/lib/x86_64-linux-gnu/libuuid.so.1    /usr/lib/
COPY --from=0 /usr/lib/clang/17/include                 /usr/lib/clang/17/include
COPY --from=0 /usr/lib/swift                            /usr/lib/swift

ENTRYPOINT ["swift"]

CMD ["--version"]
