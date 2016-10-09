FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl http://www.cpan.org/src/5.0/perl-5.24.0.tar.xz | tar xJf -

RUN cd perl-5.24.0                              \
 && ./Configure                                 \
    -Accflags='-DNO_MATHOMS -DPERL_DISABLE_PMC' \
    -Dprefix=/usr                               \
    -des                                        \
 && make                                        \
 && make install-strip

RUN ls -l /usr/bin/perl # 1847696
RUN ldd /usr/bin/perl
RUN ls -l /lib/ld-musl-x86_64.so.1
