FROM alpine:3.23 AS builder

RUN apk add --no-cache autoconf cmake g++ gfortran linux-headers make python3

# Fix GCC 15 with cstdint, see https://github.com/JuliaLang/julia/issues/58478
ENV CXXFLAGS='-include cstdint -flto' LDFLAGS='-flto' VER=1.12.6

ADD --unpack https://github.com/JuliaLang/julia/releases/download/v$VER/julia-$VER.tar.gz /

WORKDIR /julia-$VER

RUN make -j`nproc` prefix=/usr MARCH=x86-64-v3 install

RUN strip -s /usr/bin/julia \
             /usr/lib/julia/*.so

WORKDIR /usr/lib/julia

# Delete unused libraries.
RUN rm libcurl.so* libgfortran.so*

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/julia            /usr/bin/
COPY --from=0 /usr/lib/libcrypto.so.3   \
              /usr/lib/libgcc_s.so.1    \
              /usr/lib/libgfortran.so.5 \
              /usr/lib/libjulia.so.1.12 \
              /usr/lib/libquadmath.so.0 \
              /usr/lib/libssl.so.3      \
              /usr/lib/libstdc++.so.6   \
              /usr/lib/libzstd.so.1     /usr/lib/
COPY --from=0 /usr/lib/julia            /usr/lib/julia

ENTRYPOINT ["julia"]

CMD ["--version"]
