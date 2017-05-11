FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl http://www.cpan.org/src/5.0/perl-5.26.0-RC1.tar.xz | tar xJf -

RUN cd perl-5.26.0-RC1                          \
 && ./Configure                                 \
    -Accflags='-DNO_MATHOMS -DPERL_DISABLE_PMC' \
    -Dprefix=/usr                               \
    -des                                        \
 && make                                        \
 && make install-strip
