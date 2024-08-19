FROM alpine:3.20 AS builder

RUN apk add --no-cache curl python3

WORKDIR /usr/src

RUN curl -L https://github.com/dloscutoff/pip/archive/refs/tags/v1.2.0.tar.gz | tar xz --strip-components 1

FROM codegolf/lang-base

COPY --from=0 /bin/busybox /bin/chmod /bin/sh /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1      /lib/
COPY --from=0 /usr/src                      /usr/bin
COPY --from=0 /usr/bin/env /usr/bin/python3 /usr/bin/
COPY --from=0 /usr/lib/libpython3.12.so.1.0 /usr/lib/
COPY --from=0 /usr/lib/python3.12           /usr/lib/python3.12

RUN chmod +x /usr/bin/pip.py

ENTRYPOINT ["pip.py"]

CMD ["--version"]
