ARG VER=10.0.2

FROM swipl:$VER AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes binutils

RUN rm -r /usr/lib/swipl/demo \
          /usr/lib/swipl/doc

RUN strip /usr/bin/swipl

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                \
              /lib/x86_64-linux-gnu/libcrypto.so.3           \
              /lib/x86_64-linux-gnu/libgcc_s.so.1            \
              /lib/x86_64-linux-gnu/libgmp.so.10             \
              /lib/x86_64-linux-gnu/libm.so.6                \
              /lib/x86_64-linux-gnu/libpcre2-8.so.0          \
              /lib/x86_64-linux-gnu/libstdc++.so.6           \
              /lib/x86_64-linux-gnu/libtcmalloc_minimal.so.4 \
              /lib/x86_64-linux-gnu/libtinfo.so.6            \
              /lib/x86_64-linux-gnu/libz.so.1                /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                    /lib64/
COPY --from=0 /usr/bin/swipl                                 /usr/bin/
COPY --from=0 /usr/lib/locale                                /usr/lib/locale
COPY --from=0 /usr/lib/swipl                                 /usr/lib/swipl

ENTRYPOINT ["swipl"]

CMD ["--version"]
