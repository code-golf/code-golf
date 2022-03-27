FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base clang curl fts-dev

RUN curl -L https://github.com/jsoftware/jsource/archive/j902-release-b.tar.gz \
  | tar xz

ENV CC=clang j64x=j64

RUN mv jsource-* jsource \
 && cd jsource/make2     \
 && ./build_jconsole.sh  \
 && ./build_libj.sh      \
 && strip ../bin/linux/j64/jconsole ../bin/linux/j64/libj.so

COPY j.c /

RUN gcc -s -o j j.c

FROM scratch

COPY --from=0 /empty                            /dev
COPY --from=0 /lib/ld-musl-x86_64.so.1          /lib/
COPY --from=0 /empty                            /proc
COPY --from=0 /empty                            /tmp
COPY --from=0 /jsource/jlibrary/bin/profile.ijs /usr/bin/profile.ijs
COPY --from=0 /jsource/jlibrary/system          /usr/system
COPY --from=0 /j                                \
              /jsource/bin/linux/j64/jconsole   \
              /jsource/bin/linux/j64/libj.so    /usr/bin/

ENTRYPOINT ["j"]

CMD ["-v"]
