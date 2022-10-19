FROM alpine:3.16 as builder

RUN apk add --no-cache R

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1   \
              /lib/libz.so.1             /lib/
COPY --from=0 /bin/sed                   \
              /bin/mkdir                 \
              /bin/uname                 \
              /bin/rm                    \
              /bin/sh                    /bin/
COPY --from=0 /usr/bin/Rscript           \
              /usr/bin/R                 /usr/bin/
COPY --from=0 /usr/lib/R                 /usr/lib/R
COPY --from=0 /usr/lib/libopenblas.so.3  \
              /usr/lib/libpcre2-8.so.0   \
              /usr/lib/liblzma.so.5      \
              /usr/lib/libbz2.so.1       \
              /usr/lib/libicuuc.so.71    \
              /usr/lib/libicui18n.so.71  \
              /usr/lib/libgomp.so.1      \
              /usr/lib/libgfortran.so.5  \
              /usr/lib/libncursesw.so.6  \
              /usr/lib/libicudata.so.71  \
              /usr/lib/libstdc++.so.6    \
              /usr/lib/libgcc_s.so.1     \
              /usr/lib/libquadmath.so.0  \
              /usr/lib/libreadline.so.8  /usr/lib/
COPY --from=0 /etc/R                     /etc/R
RUN mkdir /root

ENTRYPOINT ["Rscript"]

CMD ["-e", "commandArgs()"]
