FROM alpine:3.10 as builder

ENV FORCE_UNSAFE_CONFIGURE 1

RUN apk add --no-cache curl gcc make musl-dev

RUN curl https://haible.de/bruno/gnu/clisp-2.49.92.tar.bz2 | tar xjf -

RUN cd clisp-2.49.92/src                                      \
 && ../configure --ignore-absence-of-libsigsegv --prefix=/usr \
 && ./makemake --prefix=/usr > Makefile                       \
 && make                                                      \
 && make install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/clisp           /usr/bin/lisp
COPY --from=0 /usr/lib/clisp-2.49.92   /usr/lib/clisp-2.49.92

ENTRYPOINT ["lisp", "--version"]
