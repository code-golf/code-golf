FROM swipl:9.2.8 AS builder

ENV LC_ALL=C.UTF-8 VER=3.6.0

RUN apt-get update \
 && apt-get install -y curl gcc make sqlite3 unzip

RUN curl -L https://get.haskellstack.org | sh

RUN curl https://mirrors.cat.pdx.edu/pakcs/download/pakcs-$VER-amd64-Linux.tar.gz | tar xz

RUN mv pakcs-$VER /out

WORKDIR /out

RUN make -j`nproc`

RUN /out/bin/cypm update           \
 && /out/bin/cypm install runcurry \
 && mv /root/.cpm/bin/runcurry /runcurry

COPY /curry.c /

RUN gcc -Wall -Werror -Wextra -s -o /curry /curry.c

FROM codegolf/lang-base

ENV GNU=x86_64-linux-gnu

COPY --from=0 /bin                        /bin
COPY --from=0 /lib/$GNU/libacl.so.1       \
              /lib/$GNU/libattr.so.1      \
              /lib/$GNU/libbfd-2.40*.so   \
              /lib/$GNU/libc.so.6         \
              /lib/$GNU/libdl.so.2        \
              /lib/$GNU/libgcc_s.so.1     \
              /lib/$GNU/libgmp.so.10      \
              /lib/$GNU/libm.so.6         \
              /lib/$GNU/libpcre2-8.so.0   \
              /lib/$GNU/libpthread.so.0   \
              /lib/$GNU/libreadline.so.8  \
              /lib/$GNU/librt.so.1        \
              /lib/$GNU/libselinux.so.1   \
              /lib/$GNU/libsframe.so.0    \
              /lib/$GNU/libstdc++.so.6    \
              /lib/$GNU/libtcmalloc*.so.4 \
              /lib/$GNU/libtinfo.so.6     \
              /lib/$GNU/libutil.so.1      \
              /lib/$GNU/libz.so.1         \
              /lib/$GNU/libzstd.so.1      /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2 /lib64/
COPY --from=0 /out/pakcsrc.default        /out/
COPY --from=0 /out/bin/curry              \
              /out/bin/pakcs-fcypp        \
              /out/bin/pakcs-frontend     /out/bin/
COPY --from=0 /out/lib                    /out/lib
COPY --from=0 /out/src/pakcs              /out/src/
COPY --from=0 /curry /runcurry            \
              /usr/bin/basename           \
              /usr/bin/dirname            \
              /usr/bin/readlink           \
              /usr/bin/tty                /usr/bin/
COPY --from=0 /usr/lib/swipl              /usr/lib/swipl

ENTRYPOINT ["curry"]

CMD ["--version"]
