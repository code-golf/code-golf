FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl gmp-dev

ENV VER=3.2

RUN curl -#L https://ftp.gnu.org/gnu/gnucobol/gnucobol-$VER.tar.xz \
  | tar xJ --strip-components 1

COPY cobol.c codegen.patch /

RUN patch -p0 < codegen.patch

RUN ./configure            \
    --prefix=/usr          \
    --without-db           \
 && make -j`nproc` install \
 && strip /usr/bin/cobc    \
          /usr/lib/libcob.so.4.2.0

RUN gcc -Wall -Werror -Wextra -o /usr/bin/cobol -s cobol.c

FROM codegolf/lang-c

COPY --from=0 /usr/bin/cobc            \
              /usr/bin/cobol           /usr/bin/
COPY --from=0 /usr/include/gmp.h       \
              /usr/include/libcob.h    /usr/include/
COPY --from=0 /usr/include/libcob      /usr/include/libcob
COPY --from=0 /usr/lib/libcob.so       \
              /usr/lib/libcob.so.4     \
              /usr/lib/libcob.so.4.2.0 \
              /usr/lib/libgmp.so.10    /usr/lib/
COPY --from=0 /usr/share/gnucobol      /usr/share/gnucobol

FROM codegolf/lang-base

COPY --from=1 . /

ENTRYPOINT ["cobol"]

CMD ["--version"]
