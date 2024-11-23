FROM debian:bullseye-slim AS builder

RUN apt-get update                 \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --yes curl g++ libpcre3-dev make zlib1g-dev

ENV VER=3.0.0

# FIXING: Only this needs to target a specific release, but yeah...
RUN curl -#L https://github.com/harbour/core/archive/master.tar.gz \
  | tar xz --strip-components 1

RUN HB_BUILD_CONTRIBS=no \
    HB_BUILD_MODE=cpp    \
    HB_BUILD_STRIP=all   \
    HB_WITH_ALLEGRO=no   \
    HB_WITH_BZIP2=no     \
    HB_WITH_CAIRO=no     \
    HB_WITH_CUPS=no      \
    HB_WITH_CURL=no      \
    HB_WITH_CURSES=no    \
    HB_WITH_EXPAT=no     \
    HB_WITH_FIREBIRD=no  \
    HB_WITH_FREEIMAGE=no \
    HB_WITH_GD=no        \
    HB_WITH_GPM=no       \
    HB_WITH_GS=no        \
    HB_WITH_LIBMAGIC=no  \
    HB_WITH_MYSQL=no     \
    HB_WITH_ODBC=no      \
    HB_WITH_PGSQL=no     \
    HB_WITH_QT=no        \
    HB_WITH_SLANG=no     \
    HB_WITH_SQLITE3=no   \
    HB_WITH_X11=no       \
    make -j`nproc` install

COPY harbour.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbour -s -static harbour.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpcre.so.3    \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/libz.so.1       /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/harbour                      /usr/bin/
COPY --from=0 /usr/local/bin/harbour                \
              /usr/local/bin/hbmk2                  /usr/local/bin/
COPY --from=0 /usr/local/include/harbour            /usr/local/include/harbour

ENTRYPOINT ["harbour"]

CMD ["--version"]
