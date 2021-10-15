FROM crystallang/crystal:1.2.0-alpine as builder

# Crystal expects /dev/null to exist.
RUN mkdir /empty /mydev && mknod -m 0666 /mydev/null c 1 3

RUN strip /usr/bin/crystal

FROM scratch

ENV CRYSTAL_CACHE_DIR=/tmp PATH=/usr/bin:/bin

COPY --from=0 /bin   /bin
COPY --from=0 /mydev /dev
COPY --from=0 /lib   /lib
COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY --from=0 /usr   /usr

ENTRYPOINT ["crystal"]

CMD ["-v"]
