FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base bzip2-dev curl curl-dev gettext-dev \
                       gfortran pcre2-dev perl tzdata xz-dev zlib-dev

ENV VER=4.5.2

RUN curl -#L https://cran.r-project.org/src/base/R-4/R-$VER.tar.xz \
  | tar xJ --strip-components 1

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

RUN strip /usr/lib/R/lib/libR.so     \
          /usr/lib/R/lib/libRblas.so \
          /usr/lib/R/lib/libRlapack.so

WORKDIR /busybox

RUN curl -#L https://busybox.net/downloads/busybox-1.36.1.tar.bz2 \
  | tar xj --strip-components 1

COPY busybox-config .config

RUN make install

FROM codegolf/lang-base

COPY --from=0 /busybox/_install/bin        /bin
COPY --from=0 /lib/ld-musl-*.so.1          /lib/
COPY --from=0 /busybox/_install/usr/bin    /usr/bin
COPY --from=0 /usr/bin/printf              /usr/bin/
COPY --from=0 /usr/bin/Rscript             /usr/bin/R
COPY --from=0 /usr/lib/libbz2.so.1         \
              /usr/lib/libgcc_s.so.1       \
              /usr/lib/libgfortran.so.5    \
              /usr/lib/libgomp.so.1        \
              /usr/lib/libintl.so.8        \
              /usr/lib/liblzma.so.5        \
              /usr/lib/libpcre2-8.so.0     \
              /usr/lib/libquadmath.so.0    \
              /usr/lib/libz.so.1           \
              /usr/lib/libzstd.so.1        /usr/lib/
COPY --from=0 /usr/lib/R/bin/R             \
              /usr/lib/R/bin/Rscript       /usr/lib/R/bin/
COPY --from=0 /usr/lib/R/bin/exec/R        /usr/lib/R/bin/exec/
COPY --from=0 /usr/lib/R/etc/Renviron      \
              /usr/lib/R/etc/ldpaths       \
              /usr/lib/R/etc/repositories  /usr/lib/R/etc/
COPY --from=0 /usr/lib/R/lib/libR.so       \
              /usr/lib/R/lib/libRblas.so   \
              /usr/lib/R/lib/libRlapack.so /usr/lib/R/lib/
COPY --from=0 /usr/lib/R/library/base      /usr/lib/R/library/base
COPY --from=0 /usr/lib/R/library/compiler  /usr/lib/R/library/compiler
COPY --from=0 /usr/lib/R/library/datasets  /usr/lib/R/library/datasets
COPY --from=0 /usr/lib/R/library/grDevices /usr/lib/R/library/grDevices
COPY --from=0 /usr/lib/R/library/graphics  /usr/lib/R/library/graphics
COPY --from=0 /usr/lib/R/library/methods   /usr/lib/R/library/methods
COPY --from=0 /usr/lib/R/library/stats     /usr/lib/R/library/stats
COPY --from=0 /usr/lib/R/library/utils     /usr/lib/R/library/utils

ENTRYPOINT ["R"]

CMD ["--version"]
