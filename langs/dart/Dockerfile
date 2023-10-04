FROM dart:3.1.3-sdk AS builder

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                      /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/lib/dart/bin/dart                /usr/lib/dart/bin/
COPY --from=0 /usr/lib/dart/include                 /usr/lib/dart/include
COPY --from=0 /usr/lib/dart/lib                     /usr/lib/dart/lib

COPY --from=0 /usr/lib/dart/bin/snapshots/dartdev.dart.snapshot        \
              /usr/lib/dart/bin/snapshots/kernel-service.dart.snapshot \
              /usr/lib/dart/bin/snapshots/kernel_worker.dart.snapshot  \
              /usr/lib/dart/bin/snapshots/

COPY dart /usr/bin/

ENTRYPOINT ["dart"]

CMD ["--version"]
