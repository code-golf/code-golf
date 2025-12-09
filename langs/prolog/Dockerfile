FROM alpine:3.23 AS builder

RUN apk add --no-cache cmake curl g++ ncurses-dev ninja zlib-dev

ENV VER=10.0.0

RUN curl -#L https://www.swi-prolog.org/download/stable/src/swipl-$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /prolog

RUN cmake -G Ninja                  \
    -DBUILD_TESTING='OFF'           \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='gcc'        \
    -DCMAKE_CXX_COMPILER='g++'      \
    -DCMAKE_INSTALL_PREFIX='/usr'   \
    -DINSTALL_DOCUMENTATION='OFF'   \
    -DINSTALL_PROLOG_SRC='OFF'      \
    -DINSTALL_TESTS='OFF'           \
    -DSWIPL_PACKAGES='OFF'          \
    -DSWIPL_SHARED_LIB='OFF' ..     \
 && ninja install

RUN strip /usr/bin/swipl

COPY prolog.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/prolog -s /prolog.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/prolog           \
              /usr/bin/swipl            /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 \
              /usr/lib/libz.so.1        /usr/lib/
COPY --from=0 /usr/lib/swipl/ABI        \
              /usr/lib/swipl/boot.prc   /usr/lib/swipl/
COPY --from=0 /usr/lib/swipl/library    /usr/lib/swipl/library

ENTRYPOINT ["prolog"]

CMD ["--version"]
