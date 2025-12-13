FROM alpine:3.23 AS builder

RUN apk add --no-cache autoconf automake bash clang cmake curl libtool linux-headers make patch texinfo

ENV CC='clang' CXX='clang++' MAKEFLAGS=-j8 VER=0.1.14

WORKDIR /egel

RUN curl -#L https://github.com/egel-lang/egel/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN cp CMakeLists.static CMakeLists.txt

COPY egel.patch .

RUN patch -p0 < egel.patch

RUN curl -#L https://github.com/fmtlib/fmt/archive/refs/tags/12.1.0.tar.gz            \
  | tar xz --directory /egel/vendor/fmt --strip-components 1                          \
 && curl -#L https://github.com/unicode-org/icu/archive/refs/tags/release-78.1.tar.gz \
  | tar xz --directory /egel/vendor/icu --strip-components 1                          \
 && curl -#L https://github.com/libffi/libffi/archive/refs/tags/v3.5.2.tar.gz         \
  | tar xz --directory /egel/vendor/libffi --strip-components 1                       \
 && curl -#L https://ftpmirror.gnu.org/gnu/lightning/lightning-2.2.3.tar.gz           \
  | tar xz --directory /egel/vendor/lightning --strip-components 1

WORKDIR /egel/vendor

RUN bash makeall.sh

WORKDIR /egel/build

RUN cmake -DCMAKE_BUILD_TYPE='MinSizeRel' .. \
 && make install && strip /usr/local/bin/egel

COPY egel.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/egel -s -static /egel.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/egel       /usr/bin/
COPY --from=0 /usr/local/bin/egel /usr/local/bin/
COPY --from=0 /usr/local/lib/egel /usr/local/lib/egel

ENTRYPOINT ["egel"]

CMD ["--version"]
