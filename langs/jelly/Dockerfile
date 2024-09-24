FROM alpine:3.20 AS builder

RUN apk add --no-cache git patch py3-pip python3

RUN git clone https://github.com/DennisMitchell/jellylanguage.git

COPY /args.patch /jellylanguage/

RUN cd jellylanguage       \
 && patch -p1 < args.patch \
 && pip3 install --break-system-packages --upgrade .

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1      /lib/
COPY --from=0 /usr/bin/jelly                \
              /usr/bin/python3              /usr/bin/
COPY --from=0 /usr/lib/libpython3.12.so.1.0 /usr/lib/
COPY --from=0 /usr/lib/python3.12           /usr/lib/python3.12

ENTRYPOINT ["jelly"]
