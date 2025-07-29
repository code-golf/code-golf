FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make mpir-dev musl-dev patch

ENV VER=3.2

RUN curl -#L https://ftp.gnu.org/gnu/gnucobol/gnucobol-$VER.tar.xz \
  | tar xJ --strip-components 1

COPY codegen.patch /

RUN patch -p0 < codegen.patch

RUN ./configure               \
    --without-db              \
 && make -j`nproc` install    \
 && strip /usr/local/bin/cobc \
          /usr/local/lib/libcob.so

RUN find /usr/local/lib -not -name 'libcob.so*' -delete

COPY cobol.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/cobol -s cobol.c

FROM codegolf/lang-c

COPY --from=0 /usr/bin/cobol                                /usr/bin/
COPY --from=0 /usr/include/mpir.h                           /usr/include/
COPY --from=0 /usr/lib/libmpir.so.23                        /usr/lib/
COPY --from=0 /usr/local/bin/cobc                           /usr/local/bin/
COPY --from=0 /usr/local/include/libcob.h                   /usr/local/include/
COPY --from=0 /usr/local/include/libcob                     /usr/local/include/libcob
COPY --from=0 /usr/local/lib                                /usr/local/lib
COPY --from=0 /usr/local/share/gnucobol/config/default.conf /usr/local/share/gnucobol/config/

ENTRYPOINT ["cobol"]

CMD ["--version"]
