FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev perl

RUN curl http://www.cpan.org/src/5.0/perl-5.26.2.tar.xz | tar xJf -

RUN cd perl-5.26.2                                                        \
 && sed -i -E 's/\(array_base/(array_base say state/' regen/feature.pl    \
 && perl regen/feature.pl                                                 \
 && ./Configure                                                           \
    -Accflags='-DNO_MATHOMS -DPERL_DISABLE_PMC -DSILENT_NO_TAINT_SUPPORT' \
    -Aldflags='-static'                                                   \
    -des                                                                  \
 && make -j`nproc` miniperl                                               \
 && strip -s miniperl

FROM scratch

COPY --from=0 /perl-5.26.2/miniperl /usr/bin/perl

ENTRYPOINT ["/usr/bin/perl", "-e", "say substr $^V, 1"]
