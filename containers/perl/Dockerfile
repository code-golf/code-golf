FROM alpine:edge as builder

RUN apk add --no-cache curl gcc make musl-dev perl

RUN curl https://cpan.metacpan.org/authors/id/X/XS/XSAWYERX/perl-5.28.0.tar.gz | tar xzf -

RUN cd perl-5.28.0                                                                \
 && sed -i -E 's/\(array_base/(array_base say signatures state/' regen/feature.pl \
 && perl regen/feature.pl                                                         \
 && ./Configure                                                                   \
    -Accflags='-DNO_MATHOMS                                                       \
    -DPERL_DISABLE_PMC                                                            \
    -DPERL_HASH_USE_SBOX32_ALSO=0                                                 \
    -DSILENT_NO_TAINT_SUPPORT'                                                    \
    -Aldflags='-static'                                                           \
    -des                                                                          \
 && make -j`nproc` miniperl                                                       \
 && strip -s miniperl

FROM scratch

COPY --from=0 /perl-5.28.0/miniperl /usr/bin/perl

ENTRYPOINT ["/usr/bin/perl", "-e", "say substr $^V, 1"]
