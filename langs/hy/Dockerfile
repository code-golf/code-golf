FROM alpine:3.20 AS builder

RUN apk add --no-cache curl py-pip

RUN curl -L https://github.com/hylang/hy/archive/refs/tags/0.29.0.tar.gz | tar xz --strip-components 1

RUN pip install hy --break-system-packages

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1      \
              /lib/libz.so.1                /lib/
COPY --from=0 /usr/bin/hy                   \
              /usr/bin/python3              /usr/bin/
COPY --from=0 /usr/lib/libpython3.12.so.1.0 /usr/lib/
COPY --from=0 /usr/lib/python3.12           /usr/lib/python3.12

ENTRYPOINT ["hy"]

CMD ["-v"]
