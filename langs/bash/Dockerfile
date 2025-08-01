FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev patch

ENV URL=https://ftpmirror.gnu.org/bash VER=5.3

RUN curl -#L $URL/bash-$VER.tar.gz | tar xz --strip-components 1

RUN curl -#L $URL/bash-$VER-patches/bash53-00[1-3] | patch -p0

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
