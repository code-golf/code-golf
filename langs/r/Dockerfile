FROM alpine:3.17 as builder

RUN apk add --no-cache R

FROM codegolf/lang-base

COPY --from=0 /bin                     /bin
COPY --from=0 /etc/R                   /etc/R
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/bin                 /usr/bin
COPY --from=0 /usr/lib                 /usr/lib

ENTRYPOINT ["Rscript"]

CMD ["--version"]
