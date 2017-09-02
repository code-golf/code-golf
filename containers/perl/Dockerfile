FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl http://www.cpan.org/src/5.0/perl-5.26.0.tar.xz | tar xJf -

RUN cd perl-5.26.0                                                        \
 && ./Configure                                                           \
    -Accflags='-DNO_MATHOMS -DPERL_DISABLE_PMC -DSILENT_NO_TAINT_SUPPORT' \
    -Aldflags='-static'                                                   \
    -des                                                                  \
 && make -j`nproc` miniperl                                               \
 && strip -s miniperl

FROM scratch

COPY --from=0 /perl-5.26.0/miniperl /usr/bin/perl

ENTRYPOINT ["/usr/bin/perl", "-e", "print $^V =~ y/v//dr"]
