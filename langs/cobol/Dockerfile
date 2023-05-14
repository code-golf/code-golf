FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl gmp-dev

RUN curl https://ftp.gnu.org/gnu/gnucobol/gnucobol-3.1.2.tar.xz | tar xJ

COPY patch /

RUN cd gnucobol-3.1.2                      \
 && patch -p0 < /patch                     \
 && ./configure --prefix=/usr --without-db \
 && make -j`nproc` install                 \
 && strip /usr/bin/cobc /usr/lib/libcob.so.4.1.0

COPY cobol.c /

RUN gcc -s -o cobol cobol.c

# HACK Put in a dir so COPY will preserve the symlinks.
RUN mkdir /usr/lib/libcob \
 && mv /usr/lib/libcob.so* /usr/lib/libcob

FROM codegolf/lang-c

COPY --from=0 /cobol /usr/bin/cobc  /usr/bin/
COPY --from=0 /usr/lib/libcob       /usr/lib
COPY --from=0 /usr/lib/libgmp.so.10 /usr/lib/
COPY --from=0 /usr/include/gmp.h    \
              /usr/include/libcob.h /usr/include/
COPY --from=0 /usr/include/libcob   /usr/include/libcob
COPY --from=0 /usr/share/gnucobol   /usr/share/gnucobol

ENTRYPOINT ["cobol"]

CMD ["--version"]
