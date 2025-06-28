FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl

ENV VER=5.40.2

RUN curl -#L https://cpan.metacpan.org/authors/id/S/SH/SHAY/perl-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./Configure                       \
    -Accflags='                       \
        -DNO_LOCALE                   \
        -DNO_MATHOMS                  \
        -DPERL_DISABLE_PMC            \
        -DPERL_HASH_USE_SBOX32_ALSO=0 \
        -DPERL_USE_SAFE_PUTENV        \
        -DSILENT_NO_TAINT_SUPPORT'    \
    -des                              \
    -Darchlib=/usr/lib/perl           \
    -Dinc_version_list=none           \
    -Dprefix=/usr                     \
    -Dprivlib=/usr/lib/perl           \
    -Dsitearch=/usr/lib/perl          \
    -Dsitelib=/usr/lib/perl           \
    -Dvendorarch=/usr/lib/perl        \
 && make -j`nproc` install            \
 && strip /usr/bin/perl

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1             /lib/
COPY --from=0 /usr/bin/perl                        /usr/bin/
COPY --from=0 /usr/lib/perl/Carp.pm                \
              /usr/lib/perl/Config.pm              \
              /usr/lib/perl/DynaLoader.pm          \
              /usr/lib/perl/Exporter.pm            \
              /usr/lib/perl/XSLoader.pm            \
              /usr/lib/perl/bigint.pm              \
              /usr/lib/perl/bignum.pm              \
              /usr/lib/perl/builtin.pm             \
              /usr/lib/perl/constant.pm            \
              /usr/lib/perl/feature.pm             \
              /usr/lib/perl/integer.pm             \
              /usr/lib/perl/overload.pm            \
              /usr/lib/perl/overloading.pm         \
              /usr/lib/perl/strict.pm              \
              /usr/lib/perl/utf8.pm                \
              /usr/lib/perl/vars.pm                \
              /usr/lib/perl/warnings.pm            /usr/lib/perl/
COPY --from=0 /usr/lib/perl/Exporter/Heavy.pm      /usr/lib/perl/Exporter/
COPY --from=0 /usr/lib/perl/File/Glob.pm           /usr/lib/perl/File/
COPY --from=0 /usr/lib/perl/List/Util.pm           /usr/lib/perl/List/
COPY --from=0 /usr/lib/perl/Math/BigFloat.pm       \
              /usr/lib/perl/Math/BigInt.pm         /usr/lib/perl/Math/
COPY --from=0 /usr/lib/perl/Math/BigInt/Calc.pm    \
              /usr/lib/perl/Math/BigInt/Lib.pm     /usr/lib/perl/Math/BigInt/
COPY --from=0 /usr/lib/perl/Scalar/Util.pm         /usr/lib/perl/Scalar/
COPY --from=0 /usr/lib/perl/auto/File/Glob/Glob.so /usr/lib/perl/auto/File/Glob/
COPY --from=0 /usr/lib/perl/auto/List/Util/Util.so /usr/lib/perl/auto/List/Util/
COPY --from=0 /usr/lib/perl/warnings/register.pm   /usr/lib/perl/warnings/

ENTRYPOINT ["perl"]

CMD ["--version"]
