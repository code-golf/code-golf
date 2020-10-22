FROM alpine:3.12 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl patch

RUN curl https://ftp.gnu.org/gnu/bash/bash-5.0.tar.gz | tar xzf - && cd bash-5.0 \
 && curl https://ftp.gnu.org/gnu/bash/bash-5.0-patches/bash50-0[01-18] | patch -p0

RUN cd bash-5.0                \
 && ./configure                \
    --disable-command-timing   \
    --disable-debugger         \
    --disable-directory-stack  \
    --disable-help-builtin     \
    --disable-job-control      \
    --disable-net-redirections \
    --disable-progcomp         \
    --disable-select           \
    --enable-static-link       \
    --without-bash-malloc      \
 && make                       \
 && strip bash

FROM scratch

COPY --from=0 /empty         /proc
COPY --from=0 /empty         /tmp
COPY --from=0 /bash-5.0/bash /usr/bin/

ENTRYPOINT ["/usr/bin/bash", "-c", "echo ${BASH_VERSION%\\([0-9]\\)-release}"]
