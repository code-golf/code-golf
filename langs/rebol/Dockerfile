FROM alpine:3.22 AS builder

RUN apk add --no-cache curl libc6-compat

ENV VER=278

RUN curl -#L https://www.rebol.com/downloads/v$VER/rebol-core-$VER-4-10.tar.gz \
  | tar xz --strip-components 1

FROM codegolf/lang-base

COPY --from=0 /lib/ld-linux-x86-64.so.2 \
              /lib/ld-musl-*.so.1       \
              /lib/libgcompat.so.0      /lib/
COPY --from=0 /lib64                    /lib64
COPY --from=0 /rebol                    /usr/bin/
COPY --from=0 /usr/lib/libobstack.so.1  \
              /usr/lib/libucontext.so.1 /usr/lib/

ENTRYPOINT ["rebol"]

CMD ["-qw", "--do", "print system/version"]
