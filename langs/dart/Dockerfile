FROM dart:3.5.4-sdk AS builder

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                               /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                \
              /lib/x86_64-linux-gnu/libdl.so.2               \
              /lib/x86_64-linux-gnu/libm.so.6                \
              /lib/x86_64-linux-gnu/libpthread.so.0          /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                    /lib64/
COPY --from=0 /usr/lib/dart/bin/dart                         /usr/bin/
COPY --from=0 /usr/lib/dart/bin/snapshots/kernel-service.*.* /usr/bin/snapshots/

COPY dartwrapper /usr/bin/

ENTRYPOINT ["dartwrapper"]

CMD ["--version"]
