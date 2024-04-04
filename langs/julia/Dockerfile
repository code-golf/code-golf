FROM alpine:3.19 AS builder

RUN apk add --no-cache autoconf build-base cmake curl gfortran linux-headers python3

ENV CXXFLAGS=-flto LDFLAGS=-flto VER=1.10.2

RUN curl -L https://github.com/JuliaLang/julia/releases/download/v$VER/julia-$VER.tar.gz \
  | tar xz

RUN make -C julia-$VER -j`nproc` prefix=/usr MARCH=x86-64-v3 install \
 && strip -s /usr/bin/julia /usr/lib/julia/*.so                      \
 && rm -r julia-$VER

# Delete unused libraries.
RUN rm /usr/lib/julia/libcurl.so* /usr/lib/julia/libgfortran.so*

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1      /lib/
COPY --from=0 /usr/bin/julia                /usr/bin/
COPY --from=0 /usr/lib/julia                /usr/lib/julia
COPY --from=0 /usr/lib/libbrotlicommon.so.1 \
              /usr/lib/libbrotlidec.so.1    \
              /usr/lib/libcares.so.2        \
              /usr/lib/libcrypto.so.3       \
              /usr/lib/libcurl.so.4         \
              /usr/lib/libgcc_s.so.1        \
              /usr/lib/libgfortran.so.5     \
              /usr/lib/libidn2.so.0         \
              /usr/lib/libjulia.so.1.10     \
              /usr/lib/libquadmath.so.0     \
              /usr/lib/libssl.so.3          \
              /usr/lib/libstdc++.so.6       \
              /usr/lib/libunistring.so.5    /usr/lib/

ENTRYPOINT ["julia"]

CMD ["-v"]
