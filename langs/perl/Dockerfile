FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl

ENV VER 5.38.0

RUN curl https://www.cpan.org/src/5.0/perl-$VER.tar.xz | tar xJ

RUN mv perl-$VER perl                                           \
 && cd perl                                                       \
 && ./Configure                                                   \
    -Accflags='                                                   \
        -DNO_LOCALE                                               \
        -DNO_MATHOMS                                              \
        -DPERL_DISABLE_PMC                                        \
        -DPERL_HASH_USE_SBOX32_ALSO=0                             \
        -DPERL_USE_SAFE_PUTENV                                    \
        -DSILENT_NO_TAINT_SUPPORT'                                \
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

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1         /lib/
COPY --from=0 /perl/perl                       /usr/bin/
COPY --from=0 /perl/lib/Carp.pm                \
              /perl/lib/Config.pm              \
              /perl/lib/DynaLoader.pm          \
              /perl/lib/Exporter.pm            \
              /perl/lib/XSLoader.pm            \
              /perl/lib/bigint.pm              \
              /perl/lib/bignum.pm              \
              /perl/lib/constant.pm            \
              /perl/lib/feature.pm             \
              /perl/lib/integer.pm             \
              /perl/lib/overload.pm            \
              /perl/lib/overloading.pm         \
              /perl/lib/strict.pm              \
              /perl/lib/utf8.pm                \
              /perl/lib/vars.pm                \
              /perl/lib/warnings.pm            /usr/lib/perl/
COPY --from=0 /perl/lib/Exporter/Heavy.pm      /usr/lib/perl/Exporter/
COPY --from=0 /perl/lib/File/Glob.pm           /usr/lib/perl/File/
COPY --from=0 /perl/lib/List/Util.pm           /usr/lib/perl/List/
COPY --from=0 /perl/lib/Math/BigFloat.pm       \
              /perl/lib/Math/BigInt.pm         /usr/lib/perl/Math/
COPY --from=0 /perl/lib/Math/BigInt/Calc.pm    \
              /perl/lib/Math/BigInt/Lib.pm     /usr/lib/perl/Math/BigInt/
COPY --from=0 /perl/lib/Scalar/Util.pm         /usr/lib/perl/Scalar/
COPY --from=0 /perl/lib/auto/File/Glob/Glob.so /usr/lib/perl/auto/File/Glob/
COPY --from=0 /perl/lib/auto/List/Util/Util.so /usr/lib/perl/auto/List/Util/
COPY --from=0 /perl/lib/warnings/register.pm   /usr/lib/perl/warnings/

ENTRYPOINT ["perl"]

CMD ["-E", "say substr $^V, 1"]
