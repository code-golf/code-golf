FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes binutils curl libatomic1 libc6-dev libedit2 libncurses6 libsqlite3-0 libz3-4

ENV VER=6.0.2

RUN curl -#L https://download.swift.org/swift-$VER-release/ubuntu2004/swift-$VER-RELEASE/swift-$VER-RELEASE-ubuntu20.04.tar.gz \
  | tar xz --strip-components 1

RUN strip /usr/bin/swift /usr/bin/swift-frontend \
 && find /usr/lib/swift -name '*.so' -exec strip {} +

RUN rm -r /usr/lib/swift/linux/x86_64

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                       \
              /lib/x86_64-linux-gnu/libdl.so.2                      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1                   \
              /lib/x86_64-linux-gnu/libm.so.6                       \
              /lib/x86_64-linux-gnu/libncurses.so.6                 \
              /lib/x86_64-linux-gnu/libpthread.so.0                 \
              /lib/x86_64-linux-gnu/librt.so.1                      \
              /lib/x86_64-linux-gnu/libtinfo.so.6                   \
              /lib/x86_64-linux-gnu/libutil.so.1                    \
              /lib/x86_64-linux-gnu/libz.so.1                       /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                           /lib64/
COPY --from=0 /usr/bin/swift /usr/bin/swift-frontend                /usr/bin/
COPY --from=0 /usr/include                                          /usr/include
COPY --from=0 /usr/lib/x86_64-linux-gnu/libsqlite3.so.0             \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6              \
              /usr/lib/x86_64-linux-gnu/libuuid.so.1                /usr/lib/
COPY --from=0 /usr/lib/swift/host/*BasicFormat.so                   \
              /usr/lib/swift/host/*CompilerPluginMessageHandling.so \
              /usr/lib/swift/host/*Diagnostics.so                   \
              /usr/lib/swift/host/*IDEUtils.so                      \
              /usr/lib/swift/host/*Operators.so                     \
              /usr/lib/swift/host/*Parser.so                        \
              /usr/lib/swift/host/*ParserDiagnostics.so             \
              /usr/lib/swift/host/*Syntax.so                        \
              /usr/lib/swift/host/*SyntaxBuilder.so                 \
              /usr/lib/swift/host/*SyntaxMacroExpansion.so          \
              /usr/lib/swift/host/*SyntaxMacros.so                  /usr/lib/swift/host/
COPY --from=0 /usr/lib/swift/linux                                  /usr/lib/swift/linux
COPY --from=0 /usr/lib/swift/shims/AssertionReporting.h             \
              /usr/lib/swift/shims/CoreFoundationShims.h            \
              /usr/lib/swift/shims/EmbeddedShims.h                  \
              /usr/lib/swift/shims/FoundationShims.h                \
              /usr/lib/swift/shims/GlobalObjects.h                  \
              /usr/lib/swift/shims/HeapObject.h                     \
              /usr/lib/swift/shims/KeyPath.h                        \
              /usr/lib/swift/shims/LibcShims.h                      \
              /usr/lib/swift/shims/MetadataSections.h               \
              /usr/lib/swift/shims/Random.h                         \
              /usr/lib/swift/shims/RefCount.h                       \
              /usr/lib/swift/shims/Reflection.h                     \
              /usr/lib/swift/shims/RuntimeShims.h                   \
              /usr/lib/swift/shims/RuntimeStubs.h                   \
              /usr/lib/swift/shims/SwiftStdbool.h                   \
              /usr/lib/swift/shims/SwiftStddef.h                    \
              /usr/lib/swift/shims/SwiftStdint.h                    \
              /usr/lib/swift/shims/System.h                         \
              /usr/lib/swift/shims/Target.h                         \
              /usr/lib/swift/shims/ThreadLocalStorage.h             \
              /usr/lib/swift/shims/UnicodeData.h                    \
              /usr/lib/swift/shims/Visibility.h                     \
              /usr/lib/swift/shims/module.modulemap                 /usr/lib/swift/shims/

ENTRYPOINT ["swift"]

CMD ["--version"]
