FROM alpine:3.22

WORKDIR /scratch

RUN mkdir dev etc proc tmp

# /dev
RUN ln -s /proc/self/fd dev/fd     \
 && mknod -m 666 dev/null    c 1 3 \
 && mknod -m 444 dev/random  c 1 8 \
 && mknod -m 444 dev/urandom c 1 9

# /etc
RUN echo nobody:x:65534:             > etc/group \
 && echo nobody:x:65534:65534::/tmp: > etc/passwd

FROM scratch AS lang-base

COPY --from=0 /scratch /
