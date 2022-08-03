FROM swipl:8.4.2 as builder

RUN apt-get update && apt-get install -y gcc

RUN rm -r /usr/lib/swipl/demo /usr/lib/swipl/doc

COPY prolog.c /

RUN gcc -s -o prolog prolog.c

FROM codegolf/lang-base

COPY --from=0 /prolog /usr/bin/swipl                             /usr/bin/
COPY --from=0 /lib                                               /lib
COPY --from=0 /lib64                                             /lib64
COPY --from=0 /usr/lib/swipl                                     /usr/lib/swipl
COPY --from=0 /usr/lib/locale                                    /usr/lib/locale
COPY --from=0 /usr/lib/x86_64-linux-gnu/libcrypto.so.1.1         \
              /usr/lib/x86_64-linux-gnu/libtcmalloc_minimal.so.4 \
              /usr/lib/x86_64-linux-gnu/libgmp.so.10             \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6           /usr/lib/

ENTRYPOINT ["prolog"]

CMD ["--version"]
