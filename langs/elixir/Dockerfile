# FIXME Using alpine causes a 5s timeout with no net, musl DNS???
# <1s docker run --rm --net none elixir:slim
# >5s docker run --rm --net none elixir:alpine
FROM elixir:1.15.4-slim AS builder

FROM codegolf/lang-base

COPY --from=0 /bin/readlink                             \
              /bin/sed                                  \
              /bin/sh                                   /bin/
COPY --from=0 /lib/x86_64-linux-gnu/ld-2.31.so          \
              /lib/x86_64-linux-gnu/libc.so.6           \
              /lib/x86_64-linux-gnu/libdl.so.2          \
              /lib/x86_64-linux-gnu/libgcc_s.so.1       \
              /lib/x86_64-linux-gnu/libm.so.6           \
              /lib/x86_64-linux-gnu/libpthread.so.0     \
              /lib/x86_64-linux-gnu/librt.so.1          \
              /lib/x86_64-linux-gnu/libselinux.so.1     \
              /lib/x86_64-linux-gnu/libtinfo.so.6       /lib/x86_64-linux-gnu/
COPY --from=0 /lib64                                    /lib64
COPY --from=0 /usr/bin/basename                         \
              /usr/bin/dirname                          /usr/bin/
COPY --from=0 /usr/lib/locale                           /usr/lib/locale
COPY --from=0 /usr/lib/x86_64-linux-gnu/libacl.so.1     \
              /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0 \
              /usr/lib/x86_64-linux-gnu/libstdc++.so.6  /usr/lib/x86_64-linux-gnu/
COPY --from=0 /usr/local/bin                            /usr/local/bin
COPY --from=0 /usr/local/lib                            /usr/local/lib

ENTRYPOINT ["elixir"]

CMD ["--short-version"]
