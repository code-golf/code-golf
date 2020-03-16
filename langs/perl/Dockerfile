FROM alpine:3.11 as builder

RUN apk add --no-cache build-base curl perl

RUN curl http://www.cpan.org/src/5.0/perl-5.30.2.tar.xz | tar xJf -

RUN cd perl-5.30.2                                   \
 && sed -i s/error=// cflags.SH                      \
 && sed -E 's/qw\(\)/"say"/' regen/feature.pl | perl \
 && ./Configure                                      \
    -Accflags='-DNO_LOCALE                           \
    -DNO_MATHOMS                                     \
    -DPERL_DISABLE_PMC                               \
    -DPERL_HASH_USE_SBOX32_ALSO=0                    \
    -DPERL_USE_SAFE_PUTENV                           \
    -DSILENT_NO_TAINT_SUPPORT'                       \
    -Aldflags='-static'                              \
    -des                                             \
 && make -j`nproc` miniperl                          \
 && strip -s miniperl

FROM scratch

COPY --from=0 /perl-5.30.2/miniperl /usr/bin/perl

ENTRYPOINT ["/usr/bin/perl", "-e", "say substr $^V, 1"]
