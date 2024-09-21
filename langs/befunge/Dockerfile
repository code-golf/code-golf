FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl -Lo bef.zip https://catseye.tc/distfiles/befunge-93-2.25.zip

RUN unzip bef.zip \
 && mv befunge-* bef

RUN cd bef && make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /bef/bin/bef             /usr/bin/

COPY /befunge /usr/bin/

ENTRYPOINT ["befunge"]
