FROM crystallang/crystal:1.0.0-alpine as builder

# Crystal expects /dev/null to exist
RUN mkdir /empty /mydev             \
 && mknod -m 0666 /mydev/null c 1 3 \
 && rm -r /usr/share/doc

FROM scratch

COPY --from=0 /bin   /bin
COPY --from=0 /mydev /dev
COPY --from=0 /lib   /lib
COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY --from=0 /usr   /usr

ENTRYPOINT ["crystal", "-v"]
