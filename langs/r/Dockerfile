FROM alpine:3.17 as builder

RUN apk add --no-cache R

FROM codegolf/lang-base

# FIXME rm, sed, sh, uname, which are all copies of the entire busybox.
#       They were symlinks but that's how docker copy works.
#       Ideally we'd build just those parts of busybox in non-symlink mode?
COPY --from=0 /bin/rm                    \
              /bin/sed                   \
              /bin/sh                    \
              /bin/uname                 /bin/
COPY --from=0 /etc/R                     /etc/R
COPY --from=0 /lib/ld-musl-x86_64.so.1   \
              /lib/libz.so.1             /lib/
COPY --from=0 /usr/bin/R                 \
              /usr/bin/Rscript           \
              /usr/bin/which             /usr/bin/
COPY --from=0 /usr/lib/R                 /usr/lib/R
COPY --from=0 /usr/lib/libbz2.so.1       \
              /usr/lib/libgcc_s.so.1     \
              /usr/lib/libgfortran.so.5  \
              /usr/lib/libgomp.so.1      \
              /usr/lib/libicudata.so.72  \
              /usr/lib/libicui18n.so.72  \
              /usr/lib/libicuuc.so.72    \
              /usr/lib/liblzma.so.5      \
              /usr/lib/libncursesw.so.6  \
              /usr/lib/libopenblas.so.3  \
              /usr/lib/libpcre2-8.so.0   \
              /usr/lib/libquadmath.so.0  \
              /usr/lib/libreadline.so.8  \
              /usr/lib/libstdc++.so.6    /usr/lib/

ENTRYPOINT ["Rscript"]

CMD ["--version"]
