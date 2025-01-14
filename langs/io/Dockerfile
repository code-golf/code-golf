FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes cmake curl git

ENV VER=7113659

RUN curl -#L https://github.com/IoLanguage/io/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /deps

RUN git clone https://github.com/kgabis/parson.git

WORKDIR /build

RUN cmake ..                      \
    -DCMAKE_BUILD_TYPE='Release'  \
    -DCMAKE_INSTALL_PREFIX='/usr' \
 && make -j`nproc` && make install

RUN strip /usr/lib/libbasekit.so \
          /usr/lib/libiovmall.so

COPY iowrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iowrapper -s /iowrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6 \
              /lib/x86_64-linux-gnu/libm.so.6 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2     /lib64/
COPY --from=0 /usr/bin/io /usr/bin/iowrapper  /usr/bin/
COPY --from=0 /usr/lib/libbasekit.so          \
              /usr/lib/libcoroutine.so        \
              /usr/lib/libgarbagecollector.so \
              /usr/lib/libiovmall.so          /usr/lib/

ENTRYPOINT ["iowrapper"]

CMD ["--version"]
