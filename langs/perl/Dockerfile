FROM alpine:3.13 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl perl

RUN curl http://www.cpan.org/src/5.0/perl-5.32.1.tar.xz | tar xJf -

RUN cd perl-5.32.1                                            \
 && perl -pe 's/qw\(indirect/$& say/' regen/feature.pl | perl \
 && ./Configure                                               \
    -Accflags='-DNO_LOCALE                                    \
    -DNO_MATHOMS                                              \
    -DPERL_DISABLE_PMC                                        \
    -DPERL_HASH_USE_SBOX32_ALSO=0                             \
    -DPERL_USE_SAFE_PUTENV                                    \
    -DSILENT_NO_TAINT_SUPPORT'                                \
    -des                                                      \
    -Dprefix=/usr                                             \
 && make -j`nproc`                                            \
 && strip -s perl

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /perl-5.32.1/perl        /usr/bin/perl
COPY --from=0 /perl-5.32.1/lib/utf8.pm /usr/lib/perl5/5.32.1/

ENTRYPOINT ["perl", "-e", "say substr $^V, 1"]
