FROM alpine:3.13 as builder

RUN mkdir /empty

ENV FORCE_UNSAFE_CONFIGURE 1

RUN apk add --no-cache build-base curl

RUN curl https://haible.de/bruno/gnu/clisp-2.49.92.tar.bz2 | tar xjf -

RUN cd clisp-2.49.92/src                                      \
 && ../configure --ignore-absence-of-libsigsegv --prefix=/usr \
 && ./makemake --prefix=/usr > Makefile                       \
 && make                                                      \
 && make install

RUN echo nobody:x:99:99::/: > /etc/passwd

FROM scratch

COPY --from=0 /etc/passwd              /etc/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin/clisp           /usr/bin/lisp
COPY --from=0 /usr/lib/clisp-2.49.92   /usr/lib/clisp-2.49.92

ENTRYPOINT ["lisp", "--version"]
