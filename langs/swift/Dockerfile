FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes binutils curl libatomic1 libc6-dev libedit2 libncurses6 libsqlite3-0 libz3-4

ENV VER=6.0.2

RUN curl -#L https://download.swift.org/swift-$VER-release/ubuntu2004/swift-$VER-RELEASE/swift-$VER-RELEASE-ubuntu20.04.tar.gz \
  | tar xz --strip-components 1

RUN strip /usr/bin/swift /usr/bin/swift-frontend \
 && find /usr/lib/swift -name '*.so' -exec strip {} +

WORKDIR /usr/lib/swift/shims

RUN rm LibcOverlayShims.h  \
       _SwiftConcurrency.h \
       _SwiftDistributed.h \
       _SynchronizationShims.h

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                                  \
              /lib/x86_64-linux-gnu/libdl.so.2                                 \
              /lib/x86_64-linux-gnu/libgcc_s.so.1                              \
              /lib/x86_64-linux-gnu/libm.so.6                                  \
              /lib/x86_64-linux-gnu/libncurses.so.6                            \
              /lib/x86_64-linux-gnu/libpthread.so.0                            \
              /lib/x86_64-linux-gnu/librt.so.1                                 \
              /lib/x86_64-linux-gnu/libtinfo.so.6                              \
              /lib/x86_64-linux-gnu/libutil.so.1                               \
              /lib/x86_64-linux-gnu/libz.so.1                                  /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                                      /lib64/
COPY --from=0 /usr/bin/swift /usr/bin/swift-frontend                           /usr/bin/
COPY --from=0 /usr/include                                                     /usr/include
COPY --from=0 /usr/lib/x86_64-linux-gnu/libsqlite3.so.0                        \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6                         \
              /usr/lib/x86_64-linux-gnu/libuuid.so.1                           /usr/lib/
COPY --from=0 /usr/lib/swift/host/libSwiftBasicFormat.so                       \
              /usr/lib/swift/host/libSwiftCompilerPluginMessageHandling.so     \
              /usr/lib/swift/host/libSwiftDiagnostics.so                       \
              /usr/lib/swift/host/libSwiftIDEUtils.so                          \
              /usr/lib/swift/host/libSwiftOperators.so                         \
              /usr/lib/swift/host/libSwiftParser.so                            \
              /usr/lib/swift/host/libSwiftParserDiagnostics.so                 \
              /usr/lib/swift/host/libSwiftSyntax.so                            \
              /usr/lib/swift/host/libSwiftSyntaxBuilder.so                     \
              /usr/lib/swift/host/libSwiftSyntaxMacroExpansion.so              \
              /usr/lib/swift/host/libSwiftSyntaxMacros.so                      /usr/lib/swift/host/
COPY --from=0 /usr/lib/swift/linux/libBlocksRuntime.so                         \
              /usr/lib/swift/linux/libFoundation.so                            \
              /usr/lib/swift/linux/libFoundationEssentials.so                  \
              /usr/lib/swift/linux/libFoundationInternationalization.so        \
              /usr/lib/swift/linux/lib_FoundationICU.so                        \
              /usr/lib/swift/linux/libdispatch.so                              \
              /usr/lib/swift/linux/libswiftCore.so                             \
              /usr/lib/swift/linux/libswiftDispatch.so                         \
              /usr/lib/swift/linux/libswiftGlibc.so                            \
              /usr/lib/swift/linux/libswiftSwiftOnoneSupport.so                \
              /usr/lib/swift/linux/libswiftSynchronization.so                  \
              /usr/lib/swift/linux/libswift_Concurrency.so                     \
              /usr/lib/swift/linux/libswift_RegexParser.so                     \
              /usr/lib/swift/linux/libswift_StringProcessing.so                /usr/lib/swift/linux/
COPY --from=0 /usr/lib/swift/linux/Swift.swiftmodule/*.swiftmodule             /usr/lib/swift/linux/Swift.swiftmodule/
COPY --from=0 /usr/lib/swift/linux/SwiftOnoneSupport.swiftmodule/*.swiftmodule /usr/lib/swift/linux/SwiftOnoneSupport.swiftmodule/
COPY --from=0 /usr/lib/swift/linux/_Concurrency.swiftmodule/*.swiftmodule      /usr/lib/swift/linux/_Concurrency.swiftmodule/
COPY --from=0 /usr/lib/swift/linux/_RegexParser.swiftmodule/*.swiftmodule      /usr/lib/swift/linux/_RegexParser.swiftmodule/
COPY --from=0 /usr/lib/swift/linux/_StringProcessing.swiftmodule/*.swiftmodule /usr/lib/swift/linux/_StringProcessing.swiftmodule/
COPY --from=0 /usr/lib/swift/shims                                             /usr/lib/swift/shims

ENTRYPOINT ["swift"]

CMD ["--version"]
