FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache ldc gcc musl-dev

# Delete some libraries
RUN rm -r                           \
    /usr/lib/libdruntime-ldc-debug* \
    /usr/lib/libphobos2-ldc-debug*  \
    /usr/lib/libphobos2-ldc.a       \
    /usr/lib/libdruntime-ldc.a      \
    /usr/lib/libc.a

# Symlink cleanup
RUN find usr/bin -type l ! -name cc -delete

FROM scratch

COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           \
              /lib/libz.so.1.2.11      /lib/
COPY --from=0 /usr/bin                 /usr/bin
COPY --from=0 /usr/lib                 /usr/lib
COPY --from=0 /usr/libexec             /usr/libexec
COPY --from=0 /usr/include             /usr/include
COPY --from=0 /etc/ldc2.conf           /etc/

ENTRYPOINT ["ldc2"]

CMD ["--version"]
