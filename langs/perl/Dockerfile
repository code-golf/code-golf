FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl perl

RUN curl http://www.cpan.org/src/5.0/perl-5.34.1.tar.xz | tar xJ

RUN mv perl-5.34.1 perl                                           \
 && cd perl                                                       \
 && perl -pe 's/default => \[/$& "say",/' regen/feature.pl | perl \
 && ./Configure                                                   \
    -Accflags='-DNO_LOCALE                                        \
    -DNO_MATHOMS                                                  \
    -DPERL_DISABLE_PMC                                            \
    -DPERL_HASH_USE_SBOX32_ALSO=0                                 \
    -DPERL_USE_SAFE_PUTENV                                        \
    -DSILENT_NO_TAINT_SUPPORT'                                    \
    -des                                                          \
    -Darchlib=/usr/lib/perl                                       \
    -Dinc_version_list=none                                       \
    -Dprefix=/usr                                                 \
    -Dprivlib=/usr/lib/perl                                       \
    -Dsitearch=/usr/lib/perl                                      \
    -Dsitelib=/usr/lib/perl                                       \
    -Dvendorarch=/usr/lib/perl                                    \
 && make -j`nproc`                                                \
 && strip -s perl

FROM scratch

COPY --from=0 /empty                   /dev
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /perl/perl               /usr/bin/
COPY --from=0 /perl/lib/Carp.pm        \
              /perl/lib/Config.pm      \
              /perl/lib/DynaLoader.pm  \
              /perl/lib/Exporter.pm    \
              /perl/lib/XSLoader.pm    \
              /perl/lib/feature.pm     \
              /perl/lib/strict.pm      \
              /perl/lib/utf8.pm        \
              /perl/lib/warnings.pm    /usr/lib/perl/
COPY --from=0 /perl/lib/File/Glob.pm   /usr/lib/perl/File/
COPY --from=0 /perl/lib/auto/File/Glob/Glob.so \
          /usr/lib/perl/auto/File/Glob/

ENTRYPOINT ["perl"]

CMD ["-e", "say substr $^V, 1"]
