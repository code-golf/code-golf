FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl -Lo b93.zip https://catseye.tc/distfiles/befunge-93-2.25.zip

RUN unzip b93.zip    \
 && mv befunge-* b93 \
 && cd b93           \
 && make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /b93/bin/bef             /usr/bin/

COPY /befunge /usr/bin/

ENTRYPOINT ["befunge"]
