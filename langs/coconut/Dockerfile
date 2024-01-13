FROM alpine:3.19 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev python3 gcc python3-dev py3-pip

ENV VER 3.0.4

RUN pip install coconut==$VER --break-system-packages

RUN cd /usr/bin && ls -d

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /lib/libc.musl-x86_64.so.1 /lib/
COPY --from=0 /lib/libz.so.1 /lib/
COPY --from=0 /usr/bin/      /usr/bin/
COPY --from=0 /usr/lib/      /usr/lib/
COPY --from=0 /usr/local/bin/      /usr/local/bin/

ENTRYPOINT ["coconut"]

CMD ["-v"]
