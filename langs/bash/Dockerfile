FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl ncurses-dev

WORKDIR /bash

RUN curl -#L https://ftp.gnu.org/gnu/bash/bash-5.2.tar.gz \
  | tar xz --strip-components 1

RUN curl https://ftp.gnu.org/gnu/bash/bash-5.2-patches/bash52-0[01-37] | patch -p0

RUN ./configure                \
    --disable-command-timing   \
    --disable-debugger         \
    --disable-directory-stack  \
    --disable-help-builtin     \
    --disable-net-redirections \
    --disable-progcomp         \
    --disable-select           \
    --prefix=/usr              \
    --without-bash-malloc      \
 && make install               \
 && strip /usr/bin/bash

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/bash             /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 /usr/lib/

ENTRYPOINT ["bash"]

CMD ["--version"]
