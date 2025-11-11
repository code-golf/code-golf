FROM alpine:3.22 AS builder

RUN apk add --no-cache autoconf automake bash cmake curl g++ libtool linux-headers make texinfo

ENV MAKEFLAGS=-j8 VER=0.1.14

WORKDIR /egel

RUN curl -#L https://github.com/egel-lang/egel/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN sed -i 's/#s/s/g' CMakeLists.static \
 && cp CMakeLists.static CMakeLists.txt \
 && sed -i 's/uin/size_/g' src/eval.hpp

RUN curl -#L https://github.com/fmtlib/fmt/archive/refs/heads/master.tar.gz      \
  | tar xz --directory /egel/vendor/fmt --strip-components 1                     \
 && curl -#L https://github.com/unicode-org/icu/archive/refs/heads/master.tar.gz \
  | tar xz --directory /egel/vendor/icu --strip-components 1                     \
 && curl -#L https://github.com/libffi/libffi/archive/refs/heads/master.tar.gz   \
  | tar xz --directory /egel/vendor/libffi --strip-components 1                  \
 && curl -#L https://ftpmirror.gnu.org/gnu/lightning/lightning-2.2.3.tar.gz      \
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
