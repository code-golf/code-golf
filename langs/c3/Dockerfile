FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl dash gcc libxml2-dev

ENV VER=0.7.7

RUN curl -#L https://github.com/c3lang/c3c/releases/download/v$VER/c3-linux.tar.gz \
  | tar xz --directory /usr/local/bin --strip-components 1

COPY c3.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/c3 -s c3.c

FROM codegolf/lang-base

COPY --from=0 /bin/dash                                         /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libbrotlicommon.so.1        \
              /lib/x86_64-linux-gnu/libbrotlidec.so.1           \
              /lib/x86_64-linux-gnu/libc.so.6                   \
              /lib/x86_64-linux-gnu/libcom_err.so.2             \
              /lib/x86_64-linux-gnu/libcrypto.so.3              \
              /lib/x86_64-linux-gnu/libcurl.so.4                \
              /lib/x86_64-linux-gnu/libffi.so.8                 \
              /lib/x86_64-linux-gnu/libgcc_s.so.1               \
              /lib/x86_64-linux-gnu/libgmp.so.10                \
              /lib/x86_64-linux-gnu/libgnutls.so.30             \
              /lib/x86_64-linux-gnu/libgssapi_krb5.so.2         \
              /lib/x86_64-linux-gnu/libhogweed.so.6             \
              /lib/x86_64-linux-gnu/libidn2.so.0                \
              /lib/x86_64-linux-gnu/libk5crypto.so.3            \
              /lib/x86_64-linux-gnu/libkeyutils.so.1            \
              /lib/x86_64-linux-gnu/libkrb5.so.3                \
              /lib/x86_64-linux-gnu/libkrb5support.so.0         \
              /lib/x86_64-linux-gnu/liblber.so.2                \
              /lib/x86_64-linux-gnu/libldap.so.2                \
              /lib/x86_64-linux-gnu/liblzma.so.5                \
              /lib/x86_64-linux-gnu/libm.so.6                   \
              /lib/x86_64-linux-gnu/libmvec.so.1                \
              /lib/x86_64-linux-gnu/libnettle.so.8              \
              /lib/x86_64-linux-gnu/libnghttp2.so.14            \
              /lib/x86_64-linux-gnu/libnghttp3.so.9             \
              /lib/x86_64-linux-gnu/libp11-kit.so.0             \
              /lib/x86_64-linux-gnu/libpsl.so.5                 \
              /lib/x86_64-linux-gnu/libresolv.so.2              \
              /lib/x86_64-linux-gnu/librtmp.so.1                \
              /lib/x86_64-linux-gnu/libsasl2.so.2               \
              /lib/x86_64-linux-gnu/libssh2.so.1                \
              /lib/x86_64-linux-gnu/libssl.so.3                 \
              /lib/x86_64-linux-gnu/libstdc++.so.6              \
              /lib/x86_64-linux-gnu/libtasn1.so.6               \
              /lib/x86_64-linux-gnu/libtinfo.so.6               \
              /lib/x86_64-linux-gnu/libunistring.so.5           \
              /lib/x86_64-linux-gnu/libxml2.so.2                \
              /lib/x86_64-linux-gnu/libz.so.1                   \
              /lib/x86_64-linux-gnu/libzstd.so.1                /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                       /lib64/
COPY --from=0 /usr/bin/c3 /usr/bin/cc /usr/bin/ld               /usr/bin/
COPY --from=0 /usr/lib/gcc/x86_64-linux-gnu/14/crtbeginS.o      \
              /usr/lib/gcc/x86_64-linux-gnu/14/crtendS.o        \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc.a         \
              /usr/lib/gcc/x86_64-linux-gnu/14/libgcc_s.so      \
              /usr/lib/gcc/x86_64-linux-gnu/14/liblto_plugin.so /usr/lib/gcc/x86_64-linux-gnu/14/
COPY --from=0 /usr/lib/x86_64-linux-gnu/Scrt1.o                 \
              /usr/lib/x86_64-linux-gnu/crti.o                  \
              /usr/lib/x86_64-linux-gnu/crtn.o                  \
              /usr/lib/x86_64-linux-gnu/libbfd-2.44-system.so   \
              /usr/lib/x86_64-linux-gnu/libc.so                 \
              /usr/lib/x86_64-linux-gnu/libc_nonshared.a        \
              /usr/lib/x86_64-linux-gnu/libctf.so.0             \
              /usr/lib/x86_64-linux-gnu/libdl.a                 \
              /usr/lib/x86_64-linux-gnu/libjansson.so.4         \
              /usr/lib/x86_64-linux-gnu/libm.so                 \
              /usr/lib/x86_64-linux-gnu/libpthread.a            \
              /usr/lib/x86_64-linux-gnu/libsframe.so.1          /usr/lib/x86_64-linux-gnu/
COPY --from=0 /usr/local/bin/c3c                                /usr/local/bin/
COPY --from=0 /usr/local/bin/lib/std                            /usr/local/bin/lib/std

ENTRYPOINT ["c3"]

CMD ["--version"]
