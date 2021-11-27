FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl patch

RUN curl https://ftp.gnu.org/gnu/bash/bash-5.1.tar.gz | tar xz && cd bash-5.1 \
 && curl https://ftp.gnu.org/gnu/bash/bash-5.1-patches/bash51-0[01-12] | patch -p0

# FIXME bash51-10 patch breaks --disable-job-control.
RUN cd bash-5.1                \
 && ./configure                \
    --disable-command-timing   \
    --disable-debugger         \
    --disable-directory-stack  \
    --disable-help-builtin     \
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
COPY --from=0 /bash-5.1/bash /usr/bin/

ENTRYPOINT ["bash"]

CMD ["-c", "echo ${BASH_VERSION%\\([0-9]\\)-release}"]
