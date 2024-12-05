FROM dart:3.5.4-sdk AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes gcc

COPY dartwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/dartwrapper -s -static /dartwrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                          \
              /lib/x86_64-linux-gnu/libdl.so.2                         \
              /lib/x86_64-linux-gnu/libm.so.6                          \
              /lib/x86_64-linux-gnu/libpthread.so.0                    /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                              /lib64/
COPY --from=0 /usr/bin/dartwrapper /usr/lib/dart/bin/dart              /usr/bin/
COPY --from=0 /usr/lib/dart/bin/snapshots/dartdev.dart.snapshot        \
              /usr/lib/dart/bin/snapshots/kernel-service.dart.snapshot \
              /usr/lib/dart/bin/snapshots/kernel_worker.dart.snapshot  /usr/bin/snapshots/
COPY --from=0 /usr/lib/dart/include                                    /usr/include
COPY --from=0 /usr/lib/dart/lib                                        /usr/lib

ENTRYPOINT ["dart"]

CMD ["--version"]
