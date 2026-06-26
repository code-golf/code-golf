FROM alpine:3.24 AS builder

RUN apk add --no-cache bzip2-dev curl-dev g++ gettext-dev gfortran make \
                       pcre2-dev perl tzdata xz-dev zlib-dev

ENV BUSYBOX_VER=1.37.0 VER=4.6.1

ADD --unpack https://cran.r-project.org/src/base/R-4/R-$VER.tar.xz /

WORKDIR /R-$VER

RUN ./configure                      \
    --disable-byte-compiled-packages \
    --disable-java                   \
    --enable-R-shlib                 \
    --enable-fast-install            \
    --enable-lto                     \
    --prefix=/usr                    \
    --without-readline               \
    --without-recommended-packages   \
    --without-x                      \
 && make -j`nproc` && make install

RUN strip /usr/lib/R/lib/libR*.so

# R needs some UNIX utils, build our own BusyBox with symlinks to save space.
ADD --unpack https://busybox.net/downloads/busybox-$BUSYBOX_VER.tar.bz2 /

WORKDIR /busybox-$BUSYBOX_VER

COPY busybox-config .config

RUN make install

FROM codegolf/lang-base

COPY --from=0 /busybox-*/_install            /
COPY --from=0 /lib/ld-musl-*.so.1            /lib/
COPY --from=0 /usr/bin/Rscript               /usr/bin/R
COPY --from=0 /usr/lib/libbz2.so.1           \
              /usr/lib/libgcc_s.so.1         \
              /usr/lib/libgfortran.so.5      \
              /usr/lib/libgomp.so.1          \
              /usr/lib/libintl.so.8          \
              /usr/lib/liblzma.so.5          \
              /usr/lib/libpcre2-8.so.0       \
              /usr/lib/libquadmath.so.0      \
              /usr/lib/libz.so.1             \
              /usr/lib/libzstd.so.1          /usr/lib/
COPY --parents                               \
     --from=0 /usr/lib/R/./bin/R             \
              /usr/lib/R/./bin/Rscript       \
              /usr/lib/R/./bin/exec/R        \
              /usr/lib/R/./etc/Renviron      \
              /usr/lib/R/./etc/ldpaths       \
              /usr/lib/R/./etc/repositories  \
              /usr/lib/R/./lib/libR.so       \
              /usr/lib/R/./lib/libRblas.so   \
              /usr/lib/R/./lib/libRlapack.so \
              /usr/lib/R/./library/base      \
              /usr/lib/R/./library/compiler  \
              /usr/lib/R/./library/datasets  \
              /usr/lib/R/./library/grDevices \
              /usr/lib/R/./library/graphics  \
              /usr/lib/R/./library/methods   \
              /usr/lib/R/./library/stats     \
              /usr/lib/R/./library/utils     /usr/lib/R/

ENTRYPOINT ["R"]

CMD ["--version"]
