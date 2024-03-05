FROM swipl:9.2.2 AS builder

RUN apt-get update && apt-get install -y gcc

RUN rm -r /usr/lib/swipl/demo /usr/lib/swipl/doc

COPY prolog.c /

RUN gcc -s -o prolog prolog.c

FROM codegolf/lang-base

COPY --from=0 /lib64/ld-linux-x86-64.so.2                        /lib64/
COPY --from=0 /prolog /usr/bin/swipl                             /usr/bin/
COPY --from=0 /usr/lib/swipl                                     /usr/lib/swipl
COPY --from=0 /usr/lib/locale                                    /usr/lib/locale
COPY --from=0 /usr/lib/x86_64-linux-gnu/libc.so.6                \
              /usr/lib/x86_64-linux-gnu/libcrypto.so.3           \
              /usr/lib/x86_64-linux-gnu/libgcc_s.so.1            \
              /usr/lib/x86_64-linux-gnu/libgmp.so.10             \
              /usr/lib/x86_64-linux-gnu/libm.so.6                \
              /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0          \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6           \
              /usr/lib/x86_64-linux-gnu/libtcmalloc_minimal.so.4 \
              /usr/lib/x86_64-linux-gnu/libtinfo.so.6            \
              /usr/lib/x86_64-linux-gnu/libz.so.1                /usr/lib/

ENTRYPOINT ["prolog"]

CMD ["--version"]
