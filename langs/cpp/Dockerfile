FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base clang

WORKDIR /runner

COPY cpp.c /runner/

RUN gcc -s -o cpp cpp.c

FROM scratch

COPY --from=0 /usr/bin/clang              \
              /usr/bin/ld                 \
              /runner/cpp                 /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    \
              /lib/libz.so.1              /lib/
COPY --from=0 /usr/include/               /usr/include/
COPY --from=0 /usr/lib/crt1.o             \
              /usr/lib/crti.o             \
              /usr/lib/crtn.o             \
              /usr/lib/libatomic.so       \
              /usr/lib/libbfd-2.37.so     \
              /usr/lib/libc.so            \
              /usr/lib/libclang-cpp.so.12 \
              /usr/lib/libctf.so.0        \
              /usr/lib/libffi.so.8        \
              /usr/lib/libgcc_s.so        \
              /usr/lib/libgcc_s.so.1      \
              /usr/lib/libLLVM-12.so      \
              /usr/lib/liblzma.so.5       \
              /usr/lib/libstdc++.so       \
              /usr/lib/libstdc++.so.6     \
              /usr/lib/libxml2.so.2       \
              /usr/lib/Scrt1.o            /usr/lib/
COPY --from=0 /usr/lib/clang/             /usr/lib/clang/
COPY --from=0 /usr/lib/gcc/               /usr/lib/gcc/
COPY --from=0 /empty                      /tmp
COPY --from=0 /empty                      /proc

ENTRYPOINT ["/usr/bin/cpp"]

CMD ["--version"]
