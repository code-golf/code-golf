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
          /usr/lib/libcob.so

RUN gcc -Wall -Werror -Wextra -o /usr/bin/cobol -s cobol.c

# Docker can only copy symlinks by copying the whole folder.
# See https://github.com/moby/moby/issues/40449
WORKDIR /usr/lib
RUN mkdir /usr-lib && cp -P libcob.so* libgmp.so* /usr-lib

FROM codegolf/lang-c

COPY --from=0 /usr/bin/cobc         \
              /usr/bin/cobol        /usr/bin/
COPY --from=0 /usr/include/gmp.h    \
              /usr/include/libcob.h /usr/include/
COPY --from=0 /usr/include/libcob   /usr/include/libcob
COPY --from=0 /usr-lib              /usr/lib
COPY --from=0 /usr/share/gnucobol   /usr/share/gnucobol

ENTRYPOINT ["cobol"]

CMD ["--version"]
