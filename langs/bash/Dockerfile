FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl patch

RUN curl https://ftp.gnu.org/gnu/bash/bash-5.2.tar.gz | tar xz && cd bash-5.2 \
 && curl https://ftp.gnu.org/gnu/bash/bash-5.2-patches/bash52-00[1-9] | patch -p0

COPY strtoimax.patch /

RUN mv bash-* bash             \
 && patch -p0 <strtoimax.patch \
 && cd bash                    \
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

FROM codegolf/lang-base

COPY --from=0 /bash/bash /usr/bin/

ENTRYPOINT ["bash"]

CMD ["-c", "echo ${BASH_VERSION%\\([0-9]\\)-release}"]
