FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes cmake curl git

ENV CFLAGS='-Wno-incompatible-pointer-types' VER=8db0702

RUN curl -#L https://github.com/IoLanguage/io/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /deps

RUN git clone https://github.com/kgabis/parson

WORKDIR /build

RUN cmake ..                        \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
 && make -j`nproc` && make install

RUN strip /usr/local/lib/libbasekit.so \
          /usr/local/lib/libiovmall.so

COPY io.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/io -s /io.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libm.so.6       /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/io                           /usr/bin/
COPY --from=0 /usr/local/lib/libbasekit.so          \
              /usr/local/lib/libcoroutine.so        \
              /usr/local/lib/libgarbagecollector.so \
              /usr/local/lib/libiovmall.so          /usr/lib/
COPY --from=0 /usr/local/bin/io                     /usr/local/bin/

ENTRYPOINT ["io"]
