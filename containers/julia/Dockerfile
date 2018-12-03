FROM alpine:edge as builder

RUN apk add --no-cache cmake curl g++ gfortran libexecinfo-dev \
    linux-headers m4 make musl-dev patch perl python3 tar xz

RUN curl -L https://github.com/JuliaLang/julia/releases/download/v1.0.2/julia-1.0.2.tar.gz | tar xzf -

RUN make -C julia prefix=/usr MARCH=x86-64 install \
 && strip -s /usr/bin/julia /usr/lib/julia/*.so    \
 && rm -r julia

COPY julia.c /

RUN gcc -s -o run-julia julia.c

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1                 /lib/
COPY --from=0 /run-julia /usr/bin/julia                /usr/bin/
COPY --from=0 /usr/lib/julia/sys.so                    /usr/lib/julia/
COPY --from=0 /usr/lib/julia/libamd.so                 \
              /usr/lib/julia/libcamd.so                \
              /usr/lib/julia/libccolamd.so             \
              /usr/lib/julia/libcholmod.so             \
              /usr/lib/julia/libcolamd.so              \
              /usr/lib/julia/libdSFMT.so               \
              /usr/lib/julia/libgmp.so                 \
              /usr/lib/julia/libgmp.so.10              \
              /usr/lib/julia/libLLVM-6.0.so            \
              /usr/lib/julia/libmpfr.so                \
              /usr/lib/julia/libopenblas64_.so         \
              /usr/lib/julia/libopenblas64_.so.0       \
              /usr/lib/julia/libpcre2-8.so             \
              /usr/lib/julia/libsuitesparse_wrapper.so \
              /usr/lib/julia/libsuitesparseconfig.so   \
              /usr/lib/libexecinfo.so.1                \
              /usr/lib/libgcc_s.so.1                   \
              /usr/lib/libgfortran.so.5                \
              /usr/lib/libjulia.so.1                   \
              /usr/lib/libquadmath.so.0                \
              /usr/lib/libstdc++.so.6                  /usr/lib/

ENTRYPOINT ["/usr/bin/julia", "-e", "println(VERSION)"]
