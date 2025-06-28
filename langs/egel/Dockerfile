FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes binutils-dev clang cmake curl libffi-dev libfmt-dev libiberty-dev libicu-dev ninja-build pkg-config zlib1g-dev

ENV VER=0.1.14
   
WORKDIR /lightning

RUN curl -#L https://ftp.gnu.org/gnu/lightning/lightning-2.2.3.tar.gz \
  | tar xz --strip-components 1

RUN ./configure && make install

WORKDIR /egel

RUN curl -#L https://github.com/egel-lang/egel/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /egel/build

RUN cmake -G Ninja .. && ninja install

COPY egel.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/egelwrapper -s /egel.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libbfd-2.40-*.so     \
              /lib/x86_64-linux-gnu/libc.so.6            \
              /lib/x86_64-linux-gnu/libffi.so.8          \
              /lib/x86_64-linux-gnu/libfmt.so.9          \
              /lib/x86_64-linux-gnu/libgcc_s.so.1        \
              /lib/x86_64-linux-gnu/libicudata.so.72     \
              /lib/x86_64-linux-gnu/libicui18n.so.72     \
              /lib/x86_64-linux-gnu/libicuio.so.72       \
              /lib/x86_64-linux-gnu/libicuuc.so.72       \
              /lib/x86_64-linux-gnu/libm.so.6            \
              /lib/x86_64-linux-gnu/libopcodes-2.40-*.so \
              /lib/x86_64-linux-gnu/libsframe.so.0       \
              /lib/x86_64-linux-gnu/libstdc++.so.6       \
              /lib/x86_64-linux-gnu/libz.so.1            \
              /lib/x86_64-linux-gnu/libzstd.so.1         /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                /lib64/
COPY --from=0 /usr/bin/egelwrapper /usr/local/bin/egel   /usr/bin/
COPY --from=0 /usr/local/lib/liblightning.so.2           /usr/lib/
COPY --from=0 /usr/local/lib/egel                        /usr/lib/egel

ENTRYPOINT ["egelwrapper"]

CMD ["--version"]
