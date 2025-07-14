FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl

WORKDIR /lang

RUN curl https://ftp.gnu.org/gnu/bash/bash-5.3.tar.gz \
  | tar xz --strip-components 1

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
    --without-bash-malloc            \
 && make                             \
 && strip bash

FROM codegolf/lang-base

COPY --from=0 /lang/bash /usr/bin/

ENTRYPOINT ["bash"]

CMD ["--version"]
