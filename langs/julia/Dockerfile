FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache cmake curl g++ gfortran libexecinfo-dev \
    linux-headers m4 make musl-dev patch perl python3 tar xz

ENV CXXFLAGS=-flto LDFLAGS=-flto VER=1.6.4

RUN curl -L https://github.com/JuliaLang/julia/releases/download/v$VER/julia-$VER.tar.gz \
  | tar xz

RUN make -C julia-$VER -j`nproc` prefix=/usr MARCH=x86-64 install \
 && strip -s /usr/bin/julia /usr/lib/julia/*.so                   \
 && rm -r julia-$VER

COPY julia.c /

RUN gcc -s -o run-julia julia.c

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1      /lib/
COPY --from=0 /empty                        /proc
COPY --from=0 /empty                        /tmp
COPY --from=0 /run-julia /usr/bin/julia     /usr/bin/
COPY --from=0 /usr/lib/julia                /usr/lib/julia/
COPY --from=0 /usr/lib/libcurl.so.4         /usr/lib/libcurl.so
COPY --from=0 /usr/lib/libbrotlicommon.so.1 \
              /usr/lib/libbrotlidec.so.1    \
              /usr/lib/libcrypto.so.1.1     \
              /usr/lib/libexecinfo.so.1     \
              /usr/lib/libgcc_s.so.1        \
              /usr/lib/libgfortran.so.5     \
              /usr/lib/libjulia.so.1        \
              /usr/lib/libquadmath.so.0     \
              /usr/lib/libssl.so.1.1        \
              /usr/lib/libstdc++.so.6       /usr/lib/

ENTRYPOINT ["/usr/bin/run-julia"]

CMD ["-v"]
