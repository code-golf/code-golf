FROM alpine:3.19 AS builder

RUN apk add --no-cache R=4.3.1-r1 build-base curl perl

RUN curl https://busybox.net/downloads/busybox-1.36.1.tar.bz2 | tar xj \
 && mv busybox-* busybox

COPY busybox-config busybox/.config

RUN cd busybox && make install

FROM codegolf/lang-base

COPY --from=0 /busybox/_install/bin      /bin
COPY --from=0 /etc/R                     /etc/R
COPY --from=0 /lib/ld-musl-x86_64.so.1   \
              /lib/libz.so.1             /lib/
COPY --from=0 /busybox/_install/usr/bin  /usr/bin
COPY --from=0 /usr/bin/R                 \
              /usr/bin/Rscript           /usr/bin/
COPY --from=0 /usr/lib/R                 /usr/lib/R
COPY --from=0 /usr/lib/libbz2.so.1       \
              /usr/lib/libgcc_s.so.1     \
              /usr/lib/libgfortran.so.5  \
              /usr/lib/libgomp.so.1      \
              /usr/lib/libicudata.so.74  \
              /usr/lib/libicui18n.so.74  \
              /usr/lib/libicuuc.so.74    \
              /usr/lib/liblzma.so.5      \
              /usr/lib/libncursesw.so.6  \
              /usr/lib/libopenblas.so.3  \
              /usr/lib/libpcre2-8.so.0   \
              /usr/lib/libquadmath.so.0  \
              /usr/lib/libreadline.so.8  \
              /usr/lib/libstdc++.so.6    /usr/lib/

ENTRYPOINT ["Rscript"]

CMD ["--version"]
