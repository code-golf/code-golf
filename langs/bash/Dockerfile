FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl ncurses-static

ENV VER=5.2

RUN curl -#L https://ftp.gnu.org/gnu/bash/bash-$VER.tar.gz \
  | tar xz --strip-components 1

RUN curl -# https://ftp.gnu.org/gnu/bash/bash-$VER-patches/bash52-0[01-37] | patch -p0

RUN ./configure                      \
    --disable-command-timing         \
    --disable-debugger               \
    --disable-directory-stack        \
    --disable-help-builtin           \
    --disable-job-control            \
    --disable-net-redirections       \
    --disable-progcomp               \
    --disable-prompt-string-decoding \
    --disable-select                 \
    --enable-static-link             \
    --prefix=/usr                    \
    --without-bash-malloc            \
 && make install                     \
 && strip /usr/bin/bash

FROM codegolf/lang-base

COPY --from=0 /usr/bin/bash /usr/bin/

ENTRYPOINT ["bash"]

CMD ["--version"]
