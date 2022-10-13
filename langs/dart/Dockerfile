FROM dart:2.18.2-sdk as builder

FROM codegolf/lang-base

COPY --from=0 /bin                                  /bin
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/lib/dart                         /usr/lib/dart

COPY dart /usr/bin/

ENTRYPOINT ["dart"]

CMD ["--version"]
