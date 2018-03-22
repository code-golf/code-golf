FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl -L http://downloads.sourceforge.net/sourceforge/clisp/clisp-2.49.tar.bz2 | tar xjf -

RUN cd clisp-2.49/src                                         \
 && sed -i '/asm\/page\.h/d' unix.d                           \
 && ../configure --ignore-absence-of-libsigsegv --prefix=/usr \
 && ./makemake --prefix=/usr > Makefile                       \
 && make                                                      \
 && make install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/clisp           /usr/bin/lisp
COPY --from=0 /usr/lib/clisp-2.49      /usr/lib/clisp-2.49

ENTRYPOINT ["lisp", "--version"]
